package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ashtishad/instabid-wallet/lib"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	FindByCredential(ctx context.Context, req LoginRequest) (*Login, lib.APIError)
}

type AuthRepoDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewAuthRepoDB(db *sql.DB, l *slog.Logger) *AuthRepoDB {
	return &AuthRepoDB{
		db: db,
		l:  l,
	}
}

func (d *AuthRepoDB) FindByCredential(ctx context.Context, req LoginRequest) (*Login, lib.APIError) {
	var sqlQuery, value, dbField string

	// prepare query according to credential field, one of email or username
	switch ctx.Value(UserCredentialKey) {
	case UserCredentialEmail:
		value = req.Email
		dbField = UserCredentialEmail
		sqlQuery = `select user_id, username, email, hashed_pass, role, status from users where email = $1`
	case UserCredentialUsername:
		value = req.Username
		dbField = UserCredentialUsername
		sqlQuery = `select user_id, username, email, hashed_pass, role, status from users where username = $1`
	default:
		return nil, lib.BadRequestError("credential field must be one of email or username")
	}

	var l Login
	var hashedPassDB []byte
	err := d.db.QueryRowContext(ctx, sqlQuery, value).Scan(&l.UserID,
		&l.Username, &l.Email, &hashedPassDB, &l.Role, &l.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.NotFoundError(fmt.Sprintf("user with %s:%s not found", dbField, value))
		}

		d.l.ErrorContext(ctx, fmt.Sprintf("unable to query user by %s", dbField), "err", err.Error())

		return nil, lib.InternalServerError(lib.ErrUnexpectedDatabase, err)
	}

	if err = bcrypt.CompareHashAndPassword(hashedPassDB, []byte(req.Password)); err != nil {
		d.l.ErrorContext(ctx, "unable to match hashed pass", "err", err.Error())
		return nil, lib.BadRequestError("input password is wrong")
	}

	return &l, nil
}
