package app

import (
	"context"
	"net/http"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
	"github.com/ashtishad/instabid-wallet/auth-api/service"
	"github.com/ashtishad/instabid-wallet/lib/jwtutils"
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

// VerifyHandler is a function to verify token and return user claims
func (ah AuthHandlers) VerifyHandler(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	token, err := jwtutils.ParseAndValidateToken(tokenStr)
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
