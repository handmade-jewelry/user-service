package verification

import (
	"context"
	"time"

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

func (r *repository) createVerification(ctx context.Context, userID int64, token string, expiredAt time.Time) error {
	query, args, err := queryBuilder.
		Insert("verification").
		Columns("user_id", "token", "is_used", "expired_at").
		Values(userID, token, false, expiredAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}

func (r *repository) getVerification(ctx context.Context, tx pgx.Tx, token string) (*Verification, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(verificationTable).
		Where(squirrel.Eq{"token": token}).
		Suffix("FOR UPDATE").
		ToSql()
	if err != nil {
		return nil, err
	}

	var verification Verification
	err = pgxscan.Get(ctx, tx, &verification, query, args...)
	if err != nil {
		return nil, err
	}

	return &verification, nil
}

func (r *repository) markTokenUsed(ctx context.Context, tx pgx.Tx, token string) error {
	query, args, err := queryBuilder.
		Update("verification").
		Set("is_used", true).
		Where(squirrel.Eq{"token": token}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}
