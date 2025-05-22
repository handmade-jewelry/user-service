package role

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

var queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type repository struct {
	dbPool *pgxpool.Pool
}

func newRepository(dbPool *pgxpool.Pool) *repository {
	return &repository{
		dbPool: dbPool,
	}
}

func (r *repository) getRoles(ctx context.Context, userID int64) ([]Role, error) {
	sql, args, err := squirrel.
		Select("r.*").
		From("role r").
		Join("user_role ur ON ur.role_id = r.id").
		Where(squirrel.Eq{"ur.user_id": userID}).
		Where("r.deleted_at IS NULL").
		Where("ur.deleted_at IS NULL").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var roles []Role
	err = pgxscan.Select(ctx, r.dbPool, &roles, sql, args...)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
