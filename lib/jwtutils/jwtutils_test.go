package jwtutils

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"os"
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

func TestBuildVerifyURL(t *testing.T) {
	t.Cleanup(func() {
		if err := os.Unsetenv("API_HOST"); err != nil {
			t.Errorf("Failed to unset API_HOST: %v", err)
		}
		if err := os.Unsetenv("AUTH_API_PORT"); err != nil {
			t.Errorf("Failed to unset AUTH_API_PORT: %v", err)
		}
		if err := os.Unsetenv("API_SCHEME"); err != nil {
			t.Errorf("Failed to unset API_SCHEME: %v", err)
		}
	})

	tests := []struct {
		name       string
		tokenStr   string
		routeName  string
		pathUserID string
		host       string
		port       string
		scheme     string
		wantURL    string
		wantErr    error
	}{
		{
			name:       "Valid",
			tokenStr:   "token",
			routeName:  "POST:/users",
			pathUserID: "123",
			host:       "localhost",
			port:       "8080",
			scheme:     "http",
			wantURL:    "http://localhost:8080/verify?routeName=POST%3A%2Fusers&token=token&userId=123",
			wantErr:    nil,
		},
		{
			name:       "EmptyToken",
			tokenStr:   "",
			routeName:  "POST:/users",
			pathUserID: "123",
			host:       "localhost",
			port:       "8080",
			scheme:     "http",
			wantURL:    "",
			wantErr:    ErrEmptyToken,
		},
		{
			name:       "EmptyRouteName",
			tokenStr:   "token",
			routeName:  "",
			pathUserID: "123",
			host:       "localhost",
			port:       "8080",
			scheme:     "http",
			wantURL:    "",
			wantErr:    ErrEmptyRouteName,
		},
		{
			name:       "MissingHost",
			tokenStr:   "token",
			routeName:  "POST:/users",
			pathUserID: "123",
			host:       "",
			port:       "8080",
			scheme:     "http",
			wantURL:    "",
			wantErr:    ErrEmptyEnvVars,
		},
		{
			name:       "InvalidScheme",
			tokenStr:   "token",
			routeName:  "POST:/users",
			pathUserID: "123",
			host:       "localhost",
			port:       "8080",
			scheme:     "",
			wantURL:    "",
			wantErr:    ErrEmptyEnvVars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv("API_HOST", tt.host); err != nil {
				t.Fatalf("Failed to set API_HOST: %v", err)
			}
			if err := os.Setenv("AUTH_API_PORT", tt.port); err != nil {
				t.Fatalf("Failed to set AUTH_API_PORT: %v", err)
			}
			if err := os.Setenv("API_SCHEME", tt.scheme); err != nil {
				t.Fatalf("Failed to set API_SCHEME: %v", err)
			}
			gotURL, err := buildVerifyURL(tt.tokenStr, tt.routeName, tt.pathUserID)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("wanted error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if gotURL != tt.wantURL {
				t.Errorf("wanted URL %s, got %s", tt.wantURL, gotURL)
			}
		})
	}
}
