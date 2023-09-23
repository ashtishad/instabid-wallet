package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/ashtishad/instabid-wallet/user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	s service.UserService
}

func (uh *UserHandlers) NewUserHandler(c *gin.Context) {
	var newUserRequest domain.NewUserReqDTO
	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
	defer cancel()

	res, apiErr := uh.s.NewUser(timeoutCtx, newUserRequest)
	if apiErr != nil {
		c.JSON(apiErr.Code(), gin.H{
			"error": apiErr.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": &res,
	})
}
