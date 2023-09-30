package app

import (
	"context"
	"net/http"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
	"github.com/ashtishad/instabid-wallet/auth-api/service"
	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	service service.AuthService
}

func (h AuthHandlers) LoginHandler(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// ToDo: Validate request
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

	res, apiErr := h.service.Login(ctx, req)
	if apiErr != nil {
		c.JSON(apiErr.Code(), gin.H{
			"error": apiErr.Error(),
		})

		return
	}

	// ToDO: Either set cookie or send token in payload
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", res.AccessToken, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": &res.AccessToken,
		"user":        &res.Login,
	})
}
