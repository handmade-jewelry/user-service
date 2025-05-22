package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo *repository
}

func NewService(dbPool *pgxpool.Pool) *Service {
	return &Service{
		repo: newRepository(dbPool),
	}
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (*User, error) {
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user, err := s.repo.createUser(ctx, email, string(hashedPassword))
	if err != nil {
		//todo переделать
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return user, nil
}

func (s *Service) hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (s *Service) checkPassword(password, usersPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(usersPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.getUser(ctx, email)
	if err != nil {
		//todo
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	isCheck, err := s.checkPassword(password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password check failed: %v", err)
	}

	if !isCheck {
		return nil, status.Error(codes.Unauthenticated, "password is incorrect")
	}

	if !user.IsVerified {
		return nil, status.Errorf(codes.PermissionDenied, "user is not verified")
	}

	return user, nil
}

func (s *Service) MarkUserIsVerified(ctx context.Context, tx pgx.Tx, userID int64) error {
	return s.repo.markUserIsVerified(ctx, tx, userID)
}
