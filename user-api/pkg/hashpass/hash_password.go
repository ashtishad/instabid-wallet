package hashpass

import (
	"context"
	"log/slog"

	"github.com/ashtishad/instabid-wallet/lib"
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 10

// Generate hashes a given password using bcrypt, with defaultCost
func Generate(ctx context.Context, pass string, l *slog.Logger) (string, lib.APIError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), defaultCost)
	if err != nil {
		l.WarnContext(ctx, lib.ErrHashingPassword, "err", err)
		return "", lib.InternalServerError(lib.ErrUnexpected, err)
	}

	return string(hashedPassword), nil
}
