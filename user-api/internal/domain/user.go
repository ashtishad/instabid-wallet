package domain

import (
	"database/sql"
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

type Profile struct {
	FirstName string
	LastName  string
	Gender    string
	Address   sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
