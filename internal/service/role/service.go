package role

import (
	"context"
	pgError "github.com/handmade-jewelry/user-service/libs/pgutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo *repository
}

func NewService(dbPool *pgxpool.Pool) *Service {
	return &Service{
		repo: newRepository(dbPool),
	}
}

func (s *Service) GetUserRolesName(ctx context.Context, userID int64) ([]string, error) {
	roles, err := s.repo.getUserRoles(ctx, userID)
	if err != nil {
		return nil, pgError.MapPostgresError("failed to get user", err)
	}

	res := make([]string, 0, len(roles))
	for _, role := range roles {
		res = append(res, string(role.Name))
	}

	return res, nil
}

func (s *Service) SetUserRole(ctx context.Context, tx pgx.Tx, userID int64, roleName RoleName) error {
	role, err := s.repo.getRoleByName(ctx, roleName)
	if err != nil {
		return pgError.MapPostgresError("failed to get role", err)
	}

	err = s.repo.setUserRole(ctx, tx, userID, role.ID)
	if err != nil {
		return pgError.MapPostgresError("failed to set user role", err)
	}

	return nil
}
