package user_verification

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/handmade-jewelry/user-service/internal/service/user"
	"github.com/handmade-jewelry/user-service/internal/service/verification"
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

func (s *Service) VerifyUserByToken(ctx context.Context, token string) (err error) {
	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("tx rollback error: %v", rbErr)
			}
		}
	}()

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

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
