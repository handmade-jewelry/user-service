package user

import (
	"github.com/handmade-jewelry/user-service/internal/service/role"
	"github.com/handmade-jewelry/user-service/internal/service/user"
	userVerification "github.com/handmade-jewelry/user-service/internal/service/user-verification"
	"github.com/handmade-jewelry/user-service/internal/service/verification"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService             *user.Service
	roleService             *role.Service
	verificationService     *verification.Service
	userVerificationService *userVerification.Service
}

func NewUserServiceServer(
	userService *user.Service,
	roleService *role.Service,
	verificationService *verification.Service,
	userVerificationService *userVerification.Service) *UserServiceServer {
	return &UserServiceServer{
		userService:             userService,
		roleService:             roleService,
		verificationService:     verificationService,
		userVerificationService: userVerificationService,
	}
}
