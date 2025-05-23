package user

import (
	"context"
	"github.com/handmade-jewelry/user-service/logger"

	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userWithRoles, err := u.userService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		logger.Error("login failed", err)
		return nil, err
	}

	return &pb.LoginResponse{
		UserId: userWithRoles.UserID,
		Roles:  convertRolesToPb(userWithRoles.Roles),
	}, nilus
}
