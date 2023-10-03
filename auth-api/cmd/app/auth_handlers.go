package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
	"github.com/ashtishad/instabid-wallet/auth-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandlers struct {
	service service.AuthService
}

func (ah AuthHandlers) LoginHandler(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	switch {
	case req.Email != "":
		ctx = context.WithValue(ctx, domain.UserCredentialKey, domain.UserCredentialEmail)
	case req.Username != "":
		ctx = context.WithValue(ctx, domain.UserCredentialKey, domain.UserCredentialUsername)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you must provide either an email or a username, along with a password."})
		return
	}

	res, apiErr := ah.service.Login(ctx, req)
	if apiErr != nil {
		c.JSON(apiErr.Code(), gin.H{
			"error": apiErr.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": &res.AccessToken,
		"user":  &res.Login,
	})
}

var HMACSecret = []byte(os.Getenv("HMACSecret"))

// VerifyHandler is a function to verify token and return user claims
func (ah AuthHandlers) VerifyHandler(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	token, err := parseAndValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

// parseAndValidateToken parses a JWT token string and validates its signature.
// The function returns the parsed token if it's valid, and an error otherwise.
// nolint:wrapcheck
func parseAndValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return HMACSecret, nil
	})
}
