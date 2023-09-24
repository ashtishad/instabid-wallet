package app

import (
	"context"
	"net/http"

	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/ashtishad/instabid-wallet/user-api/internal/service"
	"github.com/ashtishad/instabid-wallet/user-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	s service.UserService
}

func (uh *UserHandlers) NewUserHandler(c *gin.Context) {
	var newUserRequest domain.NewUserReqDTO
	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	if gin.Mode() == gin.ReleaseMode {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, utils.TimeoutCreateUser)

		defer cancel()
	}

	res, apiErr := uh.s.NewUser(ctx, newUserRequest)
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
