package verification

import (
	"context"
	"errors"
	"time"

	"crypto/rand"
	"encoding/base64"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/handmade-jewelry/user-service/internal/service/user"
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

func (s *Service) SendVerificationLink(ctx context.Context, user *user.User) error {
	token, err := s.generateToken()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to generate verification token")
	}

	expiredAt := time.Now().UTC().Add(s.verificationTokenExp)
	err = s.repo.createVerification(ctx, user.ID, token, expiredAt)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create verification: %v", err)
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
		//todo
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "verification token not found")
		}
		return nil, err
	}

	return verification, nil
}

func (s *Service) MarkTokenUsed(ctx context.Context, tx pgx.Tx, token string) error {
	return s.repo.markTokenUsed(ctx, tx, token)
}
