package user

import (
	"context"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (s *Service) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	//todo stub
	return &pb.RefreshTokenResponse{
		Token: &pb.Token{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		},
	}, nil
}
