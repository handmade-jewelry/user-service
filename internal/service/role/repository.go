package role

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
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

func (r *repository) getUserRoles(ctx context.Context, userID int64) ([]Role, error) {
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

func (r *repository) getRoleByName(ctx context.Context, roleName RoleName) (*Role, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(roleTable).
		Where("name = ?", roleName).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rle Role
	err = pgxscan.Get(ctx, r.dbPool, &rle, query, args...)
	if err != nil {
		return nil, err
	}

	return &rle, nil
}

func (r *repository) setUserRole(ctx context.Context, tx pgx.Tx, userID, roleID int64) error {
	query, args, err := squirrel.
		Insert(userRoleTable).
		Columns("user_id", "role_id").
		Values(userID, roleID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (r *repository) listRoles(ctx context.Context) ([]Role, error) {
	query, args, err := squirrel.
		Select("*").
		From(roleTable).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var roles []Role
	err = pgxscan.Select(ctx, r.dbPool, &roles, query, args...)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
