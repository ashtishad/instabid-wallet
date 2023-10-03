package jwtutils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var HMACSecret = []byte(os.Getenv("HMACSecret"))

// ParseAndValidateToken parses a JWT token string and validates its signature.
// The function returns the parsed token if it's valid, and an error otherwise.
// nolint:wrapcheck
func ParseAndValidateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return HMACSecret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, jwt.ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, jwt.ErrTokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, err
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, err
		default:
			return nil, err
		}
	}

	return token, nil
}
