package role

import "time"

type RoleName string

const (
	AdminRoleName    = RoleName("ADMIN")
	SellerRoleName   = RoleName("SELLER")
	CustomerRoleName = RoleName("CUSTOMER")
)

const (
	roleTable     = "role"
	userRoleTable = "user_role"
)

type Role struct {
	ID        int64
	Name      RoleName
	CreatedAt time.Time
	DeletedAt *time.Time
}

type UserRole struct {
	UserID     int64
	RoleID     int64
	AssignedAt time.Time
	DeletedAt  *time.Time
}
