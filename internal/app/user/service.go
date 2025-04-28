package user

import (
	pb "github.com/handmade-jewellery/user-service/pkg/api/user-service"
)

type Service struct {
	pb.UnimplementedUserServiceServer
}

func NewService() *Service {
	return &Service{}
}
