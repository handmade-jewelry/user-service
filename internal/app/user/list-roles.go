package user

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func (u *UserServiceServer) ListRoles(ctx context.Context, _ *emptypb.Empty) (*pb.GetListRolesResponse, error) {
	roles, err := u.roleService.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetListRolesResponse{
		Roles: convertListRolesToPb(roles),
	}, nil
}
