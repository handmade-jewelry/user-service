package user

import (
	"context"
	pb "github.com/handmade-jewellery/user-service/pkg/api/user-service"
)

func (s *Service) RefreshToken(context.Context, *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	//todo stub
	return nil, nil
}
