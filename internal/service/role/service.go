package role

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	pgError "github.com/handmade-jewelry/user-service/libs/pgutils"
	"github.com/handmade-jewelry/user-service/logger"
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
		logger.ErrorWithFields("failed to get user roles", err, "user_id", userID)
		return nil, pgError.MapPostgresError("roles", err)
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
		logger.ErrorWithFields("failed to get role", err, "user_id", userID)
		return pgError.MapPostgresError("role", err)
	}

	err = s.repo.setUserRole(ctx, tx, userID, role.ID)
	if err != nil {
		logger.ErrorWithFields("failed to set role", err, "user_id", userID)
		return pgError.MapPostgresError("role", err)
	}

	return nil
}

func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	roles, err := s.repo.listRoles(ctx)
	if err != nil {
		logger.Error("failed to get list roles", err)
		return nil, pgError.MapPostgresError("roles", err)
	}

	return roles, nil
}
