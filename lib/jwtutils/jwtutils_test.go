package jwtutils

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestParseAndValidateToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenFunc func() (string, error)
		wantErr   error
	}{
		{
			name: "Valid_Token",
			tokenFunc: func() (string, error) {
				return createToken(time.Now().Add(time.Hour), "")
			},
			wantErr: nil,
		},
		{
			name:      "Invalid_Token",
			tokenFunc: func() (string, error) { return "invalid_token_string", nil },
			wantErr:   jwt.ErrTokenMalformed,
		},
		{
			name:      "Empty_Token",
			tokenFunc: func() (string, error) { return "", nil },
			wantErr:   jwt.ErrTokenMalformed,
		},
		{
			name: "Wrong_Method",
			tokenFunc: func() (string, error) {
				return createTokenWithWrongMethod()
			},
			wantErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			name: "Expired_Token",
			tokenFunc: func() (string, error) {
				return createToken(time.Now().Add(-time.Hour), "")
			},
			wantErr: jwt.ErrTokenExpired,
		},
		{
			name: "Not_Valid_Yet",
			tokenFunc: func() (string, error) {
				return createToken(time.Now().Add(time.Hour*2), "nbf")
			},
			wantErr: jwt.ErrTokenNotValidYet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, err := tt.tokenFunc()
			if err != nil {
				t.Fatalf("Error generating token: %v", err)
			}

			_, err = ParseAndValidateToken(tokenStr)

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseAndValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("ParseAndValidateToken() unexpected error = %v", err)
			}
		})
	}
}

func createToken(exp time.Time, specialClaim string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	if specialClaim == "nbf" {
		claims.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Hour))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(HMACSecret)
}

func createTokenWithWrongMethod() (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
