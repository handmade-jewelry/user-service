package user

import (
	"context"
	pb "github.com/handmade-jewellery/user-service/pkg/api/user-service"
)

func (s *Service) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {
	//todo stub
	return nil, nil
}
