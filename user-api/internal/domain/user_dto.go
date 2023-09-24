package domain

import (
	"time"
)

type UserRespDTO struct {
	UserID    string    `json:"userId" binding:"required"`
	UserName  string    `json:"userName" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Status    string    `json:"status" binding:"required"`
	Role      string    `json:"role" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
}

type NewUserReqDTO struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Status   string `json:"status" binding:"_"`
	Role     string `json:"role" binding:"_"`
}
