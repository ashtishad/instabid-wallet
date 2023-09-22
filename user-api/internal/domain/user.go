package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusDeleted  = "deleted"
)

const (
	UserRoleAdmin     = "admin"
	UserRoleUser      = "user"
	UserRoleModerator = "moderator"
	UserRoleMerchant  = "merchant"
)

type User struct {
	ID         int64
	UserID     uuid.UUID
	UserName   string
	Email      string
	Status     string
	Role       string
	HashedPass string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
