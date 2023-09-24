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

func (d *UserRepoDB) Insert(ctx context.Context, u User) (*User, lib.APIError) {
	sqlInsertUserReturnID := `INSERT INTO users (username, email, hashed_pass, status, role) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`

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

func (d *UserRepoDB) findByUUID(ctx context.Context, uuid string) (*User, lib.APIError) {
	sqlFindByUUID := `SELECT user_id, username, email, status, role, created_at, updated_at from users where user_id= $1`

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

func rollbackOnError(tx *sql.Tx, err *error, l *slog.Logger) {
	if *err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			l.Warn(lib.ErrTXRollback, "rbErr", rbErr)
		}
	}
}
