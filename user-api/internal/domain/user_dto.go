package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRespDTO struct {
	UserID    uuid.UUID `json:"userId" binding:"required,len=36"`
	UserName  string    `json:"userName" binding:"required,min=7,max=64"`
	Email     string    `json:"email" binding:"required,email,max=128"`
	Status    string    `json:"status" binding:"required,oneof=active inactive deleted"`
	Role      string    `json:"role" binding:"required,oneof=user admin moderator merchant"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
}

type NewUserReqDTO struct {
	UserName string `json:"userName" binding:"required,min=7,max=64"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin moderator merchant"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive deleted"`
}
