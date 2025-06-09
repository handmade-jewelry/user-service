package user

import (
	"context"

	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) CustomerRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := u.userService.RegisterCustomer(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Result: true,
	}, nil
}
