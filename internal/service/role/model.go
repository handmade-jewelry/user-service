package role

import "time"

type Role struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	DeletedAt *time.Time
}

type UserRole struct {
	UserID     int64
	RoleID     int64
	AssignedAt time.Time
	DeletedAt  *time.Time
}
