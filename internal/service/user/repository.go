package user

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

func (r *repository) createUser(ctx context.Context, tx pgx.Tx, email, passwordHash string) (*User, error) {
	query, args, err := queryBuilder.
		Insert(usersTable).
		Columns("email", "password").
		Values(email, passwordHash).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = pgxscan.Get(ctx, tx, user, query, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) getUser(ctx context.Context, email string) (*User, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(usersTable).
		Where(squirrel.Eq{"email": email}).
		Where(squirrel.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = pgxscan.Get(ctx, r.dbPool, &user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) markUserIsVerified(ctx context.Context, tx pgx.Tx, userID int64) error {
	query, args, err := queryBuilder.
		Update(usersTable).
		Set("is_verified", true).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}
