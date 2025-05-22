package user

import (
	"time"
)

const (
	usersTable = "users"
)

type User struct {
	ID         int64
	Email      string
	Password   string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
