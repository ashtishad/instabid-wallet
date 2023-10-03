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

const (
	queryParamToken     = "token"
	mapKeyRole          = "Role"
	queryParamRouteName = "routeName"
	queryParamUserID    = "userId"
	mapKeyUserID        = "UserID"
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

func (ah AuthHandlers) VerifyHandler(c *gin.Context) {
	tokenStr := c.Query(queryParamToken)
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

	// Extract role and userId from claims
	role, roleOk := claims[mapKeyRole].(string)
	userID, userIDOk := claims[mapKeyUserID].(string)

	if !roleOk || !userIDOk {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role or UserId not found in token"})
		return
	}

	routeName := c.Query(queryParamRouteName)

	// Check role-based permissions
	if !domain.Permissions.IsAuthorizedFor(role, routeName) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
		return
	}

	// Check userId-based permissions if the path includes userId
	pathUserID := c.Query(queryParamUserID)
	if pathUserID != "" && pathUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access resources for another user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"claims": claims})
}
