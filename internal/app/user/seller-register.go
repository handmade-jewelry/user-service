package user

import (
	"context"

	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) SellerRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := u.userService.RegisterSeller(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Result: true,
	}, nil
}
