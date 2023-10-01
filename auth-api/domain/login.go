package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Login struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type AccessTokenClaims struct {
	TokenType string
	Username  string
	UserID    string
	Email     string
	Role      string
	Status    string
	jwt.RegisteredClaims
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	expirationTime := time.Now().Add(AccessTokenDuration)

	return AccessTokenClaims{
		TokenType:        TokenTypeAccess,
		Username:         l.Username,
		UserID:           l.UserID,
		Email:            l.Email,
		Role:             l.Role,
		Status:           l.Status,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)},
	}
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Login
}
