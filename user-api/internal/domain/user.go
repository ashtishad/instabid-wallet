package domain

import (
	"time"
)

type User struct {
	ID         int64
	UserID     string
	UserName   string
	Email      string
	Status     string
	Role       string
	HashedPass string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
