package role

import (
	"context"
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
	roles, err := s.repo.getRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(roles))
	for _, role := range roles {
		res = append(res, role.Name)
	}

	return res, nil
}
