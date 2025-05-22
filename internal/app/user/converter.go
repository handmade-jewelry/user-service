package user

import pb "github.com/handmade-jewelry/user-service/pkg/api/user-service"

func convertRolesToPb(roles []string) []*pb.Role {
	pbRoles := make([]*pb.Role, 0, len(roles))
	for _, role := range roles {
		pbRoles = append(pbRoles, &pb.Role{
			Name: role,
		})
	}

	return pbRoles
}
