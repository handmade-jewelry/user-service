package user

import (
	"context"

	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	roleNames, err := u.roleService.GetUserRolesName(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserRolesResponse{
		Roles: convertRolesToPb(roleNames),
	}, nil
}
