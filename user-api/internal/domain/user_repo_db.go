package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ashtishad/instabid-wallet/lib"
)

type UserRepoDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewUserRepoDB(db *sql.DB, l *slog.Logger) *UserRepoDB {
	return &UserRepoDB{
		db: db,
		l:  l,
	}
}

// Insert adds a new user to the database and returns the inserted User object.
// The method performs a transaction with isolation level read committed,
// then checks for existing usernames and emails, returned 409 conflict error if exists,
// Returns 404 or 500 if other error occurs.
func (d *UserRepoDB) Insert(ctx context.Context, u User) (*User, lib.APIError) {
	sqlInsertUserReturnID := `INSERT INTO users (username, email, hashed_pass, status, role) 
							  VALUES ($1, $2, $3, $4, $5) RETURNING user_id`

	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		d.l.WarnContext(ctx, lib.ErrTXBegin, "err", err, "isolation", sql.LevelReadCommitted)
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	defer rollbackOnError(tx, &err, d.l)

	if apiErr := d.checkExists(ctx, u.Email, u.UserName); apiErr != nil {
		return nil, apiErr
	}

	row := tx.QueryRowContext(ctx, sqlInsertUserReturnID, u.UserName, u.Email, u.HashedPass, u.Status, u.Role)
	if err = row.Scan(&u.UserID); err != nil {
		d.l.WarnContext(ctx, lib.ErrScanRow, "err", err)
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		d.l.WarnContext(ctx, lib.ErrTXCommit, "err", err)
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	return d.findByUUID(ctx, u.UserID)
}

// InsertProfile takes uuid as string of user, and inserts profile to that specific user.
// creates user profile in a transaction with isolation level read committed, returns *Profile
// if error happens it returns 404,500
func (d *UserRepoDB) InsertProfile(ctx context.Context, uuid string, up Profile) (*Profile, lib.APIError) {
	id, apiErr := d.findIDByUUID(ctx, uuid)
	if apiErr != nil {
		return nil, apiErr
	}

	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		d.l.ErrorContext(ctx, lib.ErrTXBegin, "err", err.Error())
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	defer rollbackOnError(tx, &err, d.l)

	sqlInsertProfile := `INSERT into user_profiles (user_id, first_name, last_name, gender, address) 
						values ($1, $2, $3, $4, $5)`

	var res sql.Result
	res, err = tx.ExecContext(ctx, sqlInsertProfile, id, up.FirstName, up.LastName, up.Gender, up.Address)

	if err != nil {
		d.l.ErrorContext(ctx, lib.ErrScanRow, "err", err.Error())
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	var ra int64
	ra, err = res.RowsAffected()

	if err != nil {
		d.l.ErrorContext(ctx, "unable to get rows affected", "err", err.Error())
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	if ra != 1 {
		err = fmt.Errorf("expected to affect 1 row, affected %d", ra)
		d.l.WarnContext(ctx, err.Error())

		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		d.l.ErrorContext(ctx, lib.ErrTXCommit, "err", err.Error())
		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	return d.findProfile(ctx, id)
}

// findIDByUUID retrieves user id int64 from user uuid
// returns 404 and 500 if error happens.
func (d *UserRepoDB) findIDByUUID(ctx context.Context, userID string) (int64, lib.APIError) {
	sqlFindIDByUUID := `SELECT id from users where user_id = $1`
	var id int64

	err := d.db.QueryRowContext(ctx, sqlFindIDByUUID, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, lib.NotFoundError("user not found by uuid")
		}

		d.l.ErrorContext(ctx, "failed to find id by uuid", "err", err.Error())

		return 0, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	return id, nil
}

// findByUUID retrieves a user by their UUID from the database.
// If the user is not found, a NotFoundError is returned, Any other errors result in an InternalServerError.
func (d *UserRepoDB) findByUUID(ctx context.Context, uuid string) (*User, lib.APIError) {
	sqlFindByUUID := `SELECT user_id, username, email, status, role, created_at, updated_at 
					 from users where user_id= $1`

	var u User
	row := d.db.QueryRowContext(ctx, sqlFindByUUID, uuid)

	err := row.Scan(&u.UserID, &u.UserName, &u.Email, &u.Status, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.NotFoundError("user not found by uuid")
		}

		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	return &u, nil
}

// findProfile retrieves a user profile by their user id from the database.
// If the user profile is not found, a NotFoundError is returned, Any other errors result in an InternalServerError.
func (d *UserRepoDB) findProfile(ctx context.Context, id int64) (*Profile, lib.APIError) {
	sqlFindByUUID := `SELECT  first_name, last_name, gender, address, created_at, updated_at
					 from user_profiles where user_id= $1`

	var up Profile
	row := d.db.QueryRowContext(ctx, sqlFindByUUID, id)

	err := row.Scan(&up.FirstName, &up.LastName, &up.Gender, &up.Address, &up.CreatedAt, &up.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.NotFoundError("user profile not found by user id")
		}

		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	return &up, nil
}

// checkExists checks user email or username exists in database, if any of these exists, it will return an error
// returns nil if both fields not found.
func (d *UserRepoDB) checkExists(ctx context.Context, email, username string) lib.APIError {
	const sqlCheckExists = `SELECT
    EXISTS (SELECT 1 FROM users WHERE email = $1) AS email_exists,
    EXISTS (SELECT 1 FROM users WHERE username = $2) AS username_exists;
	`

	var emailExists, usernameExists bool
	row := d.db.QueryRowContext(ctx, sqlCheckExists, email, username)

	if err := row.Scan(&emailExists, &usernameExists); err != nil {
		d.l.ErrorContext(ctx, "failed to query database", "err", err)
		return lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	switch {
	case emailExists && usernameExists:
		return lib.ConflictError(fmt.Sprintf("user with email: %s and username: %s exists", email, username))
	case emailExists:
		return lib.ConflictError(fmt.Sprintf("user with email %s exists", email))
	case usernameExists:
		return lib.ConflictError(fmt.Sprintf("user with username %s exists", username))
	default:
		return nil
	}
}

// rollbackOnError attempts to roll back the transaction if an error is present.
// It logs a warning if the rollback itself fails.
func rollbackOnError(tx *sql.Tx, err *error, l *slog.Logger) {
	if *err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			l.Warn(lib.ErrTXRollback, "rbErr", rbErr)
		}
	}
}
