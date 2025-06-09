package user_verification

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/handmade-jewelry/user-service/internal/service/user"
	"github.com/handmade-jewelry/user-service/internal/service/verification"
	pgTx "github.com/handmade-jewelry/user-service/libs/pgutils"
	"github.com/handmade-jewelry/user-service/logger"
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
			logger.Error("failed to get verification", err)
			return status.Error(codes.Internal, "internal error")
		}

		if verification.IsUsed {
			return status.Error(codes.InvalidArgument, "verification token already used")
		}

		if verification.ExpiredAt.UTC().Before(time.Now().UTC()) {
			return status.Error(codes.InvalidArgument, "verification token expired")
		}

		err = s.verificationService.MarkTokenUsed(ctx, tx, token)
		if err != nil {
			logger.ErrorWithFields("failed to mark token as used", err, "verification", verification)
			return status.Error(codes.Internal, "internal error")
		}

		err = s.userService.MarkUserIsVerified(ctx, tx, verification.UserID)
		if err != nil {
			logger.ErrorWithFields("failed to mark user as verified", err, "verification", verification)
			return status.Error(codes.Internal, "internal error")
		}

		return nil
	})

	return err
}
