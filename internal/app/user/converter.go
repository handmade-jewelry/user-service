package user

import (
	"github.com/handmade-jewelry/user-service/internal/service/role"
	pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

func convertRolesToPb(roles []string) []*pb.Role {
	pbRoles := make([]*pb.Role, 0, len(roles))
	for _, role := range roles {
		pbRoles = append(pbRoles, &pb.Role{
			Name: role,
		})
	}

	return pbRoles
}

func convertListRolesToPb(roles []role.Role) []*pb.Role {
	pbRoles := make([]*pb.Role, 0, len(roles))
	for _, role := range roles {
		pbRoles = append(pbRoles, &pb.Role{
			Name: string(role.Name),
		})
	}

	return pbRoles
}
