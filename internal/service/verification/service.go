package verification

import (
	"context"
	pgError "github.com/handmade-jewelry/user-service/libs/pgutils"
	"time"

	"crypto/rand"
	"encoding/base64"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo                 *repository
	verificationTokenExp time.Duration
}

func NewService(dbPool *pgxpool.Pool, verificationTokenExp time.Duration) *Service {
	return &Service{
		repo:                 newRepository(dbPool),
		verificationTokenExp: verificationTokenExp,
	}
}

func (s *Service) SendVerificationLink(ctx context.Context, userID int64) error {
	token, err := s.generateToken()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to generate verification token")
	}

	expiredAt := time.Now().UTC().Add(s.verificationTokenExp)
	err = s.repo.createVerification(ctx, userID, token, expiredAt)
	if err != nil {
		return pgError.MapPostgresError("failed to create verification", err)
	}

	return s.sendEmailLetter(ctx, token)
}

func (s *Service) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Service) sendEmailLetter(ctx context.Context, token string) error {
	//todo stub
	return nil
}

func (s *Service) GetVerification(ctx context.Context, tx pgx.Tx, token string) (*Verification, error) {
	verification, err := s.repo.getVerification(ctx, tx, token)
	if err != nil {
		return nil, pgError.MapPostgresError("failed to get verification token", err)
	}

	return verification, nil
}

func (s *Service) MarkTokenUsed(ctx context.Context, tx pgx.Tx, token string) error {
	err := s.repo.markTokenUsed(ctx, tx, token)
	if err != nil {
		return pgError.MapPostgresError("failed to mark used verification", err)
	}

	return nil
}
