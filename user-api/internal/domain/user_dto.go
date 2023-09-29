package domain

import (
	"time"
)

type UserRespDTO struct {
	UserID    string    `binding:"required" json:"userId"`
	UserName  string    `binding:"required" json:"userName"`
	Email     string    `binding:"required" json:"email"`
	Status    string    `binding:"required" json:"status"`
	Role      string    `binding:"required" json:"role"`
	CreatedAt time.Time `binding:"required" json:"createdAt"`
	UpdatedAt time.Time `binding:"required" json:"updatedAt"`
}

type NewUserReqDTO struct {
	UserName string `binding:"required" json:"userName"`
	Password string `binding:"required" json:"password"`
	Email    string `binding:"required" json:"email"`

	// default value will be set if empty
	// status= "active" role = "user"
	Status string `binding:"-" json:"status"`
	Role   string `binding:"-" json:"role"`
}

type ProfileRespDTO struct {
	FirstName string    `binding:"required" json:"firstName"`
	LastName  string    `binding:"required" json:"lastName"`
	Gender    string    `binding:"required" json:"gender"`
	Address   string    `binding:"-"        json:"address,omitempty"`
	CreatedAt time.Time `binding:"-"        json:"createdAt"`
	UpdatedAt time.Time `binding:"-"        json:"updatedAt"`
}

type NewProfileReqDTO struct {
	FirstName string `binding:"required" json:"firstName"`
	LastName  string `binding:"required" json:"lastName"`
	Gender    string `binding:"required" json:"gender"`
	Address   string `binding:"-"        json:"address,omitempty"`
}
