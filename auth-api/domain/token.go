package domain

import (
	"log/slog"

	"github.com/ashtishad/instabid-wallet/lib"
	"github.com/ashtishad/instabid-wallet/lib/jwtutils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthToken struct {
	token *jwt.Token
	l     *slog.Logger
}

func NewAuthToken(claims AccessTokenClaims, l *slog.Logger) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token, l: l}
}

func (t AuthToken) NewAccessToken() (string, lib.APIError) {
	signedString, err := t.token.SignedString(jwtutils.HMACSecret)
	if err != nil {
		t.l.Error("failed signing access token", "err", err.Error())
		return "", lib.InternalServerError("cannot generate access token", err)
	}

	return signedString, nil
}
