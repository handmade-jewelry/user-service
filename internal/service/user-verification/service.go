package user_verification

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/handmade-jewelry/user-service/internal/service/user"
	"github.com/handmade-jewelry/user-service/internal/service/verification"
	pgTx "github.com/handmade-jewelry/user-service/libs/pgutils"
)

type Service struct {
	userService         *user.Service
	verificationService *verification.Service
	dbPool              *pgxpool.Pool
}

func NewService(
	userService *user.Service,
	verificationService *verification.Service,
	dbPool *pgxpool.Pool,
) *Service {
	return &Service{
		userService:         userService,
		verificationService: verificationService,
		dbPool:              dbPool,
	}
}

func (s *Service) VerifyUserByToken(ctx context.Context, token string) error {
	err := pgTx.WithTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		verification, err := s.verificationService.GetVerification(ctx, tx, token)
		if err != nil {
			return err
		}

		if verification.IsUsed {
			return status.Error(codes.InvalidArgument, "verification token already used")
		}

		if verification.ExpiredAt.UTC().Before(time.Now().UTC()) {
			return status.Error(codes.InvalidArgument, "verification token expired")
		}

		err = s.verificationService.MarkTokenUsed(ctx, tx, token)
		if err != nil {
			return status.Error(codes.Internal, "failed to mark token as used")
		}

		err = s.userService.MarkUserIsVerified(ctx, tx, verification.UserID)
		if err != nil {
			return status.Error(codes.Internal, "failed to mark user as verified")
		}

		return nil
	})

	return err
}
