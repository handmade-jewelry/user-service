package user

import (
	"context"
	"github.com/handmade-jewelry/user-service/logger"

	"github.com/handmade-jewelry/user-service/internal/service/role"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) CustomerRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := u.userService.Register(ctx, req.GetEmail(), req.GetPassword(), role.CustomerRoleName)
	if err != nil {
		logger.Error("registration failed", err)
		return nil, err
	}

	return &pb.RegisterResponse{
		Result: true,
	}, nil
}
