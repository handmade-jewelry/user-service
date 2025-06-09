package pgutils

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AlreadyExistsCode = "23505"
)

func MapPostgresError(msg string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, msg+": not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case AlreadyExistsCode:
			return status.Error(codes.AlreadyExists, msg+": already exists")
		}
	}
	return status.Error(codes.Internal, "internal error")
}
