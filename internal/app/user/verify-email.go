package user

import (
	"context"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	err := u.userVerificationService.VerifyUserByToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &pb.VerifyEmailResponse{
		Result: true,
	}, nil
}
