package user

import (
	"context"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	if !validatePassword(password) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password format")
	}

	if !validateEmail(email) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format")
	}

	user, err := u.userService.LoginUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	roles, err := u.roleService.GetUserRolesName(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		UserId: user.ID,
		Roles:  convertRolesToPb(roles),
	}, nil
}
