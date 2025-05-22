package user

import (
	"context"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *UserServiceServer) Registration(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	email, password := req.GetEmail(), req.GetPassword()

	if !validatePassword(password) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password format")
	}

	if !validateEmail(email) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format")
	}

	user, err := u.userService.CreateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	err = u.verificationService.SendVerificationLink(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegistrationResponse{
		Result: true,
	}, nil
}
