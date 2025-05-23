package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/handmade-jewelry/user-service/internal/service/role"
	"github.com/handmade-jewelry/user-service/internal/service/verification"
	"github.com/handmade-jewelry/user-service/internal/util/hasher"
	"github.com/handmade-jewelry/user-service/internal/util/validation"
	pgError "github.com/handmade-jewelry/user-service/libs/pgutils"
	pgTx "github.com/handmade-jewelry/user-service/libs/pgutils"
	"github.com/handmade-jewelry/user-service/logger"
)

type Service struct {
	roleService         *role.Service
	verificationService *verification.Service
	repo                *repository
	dbPool              *pgxpool.Pool
}

func NewService(dbPool *pgxpool.Pool, roleService *role.Service, verificationService *verification.Service) *Service {
	return &Service{
		roleService:         roleService,
		verificationService: verificationService,
		repo:                newRepository(dbPool),
		dbPool:              dbPool,
	}
}

func (s *Service) Register(ctx context.Context, email, password string, roleName role.RoleName) error {
	err := validateCredentials(email, password)
	if err != nil {
		return err
	}

	user, err := s.CreateUserWithRole(ctx, email, password, roleName)
	if err != nil {
		return err
	}

	err = s.verificationService.SendVerificationLink(ctx, user.ID)
	if err != nil {
		logger.Error("failed to send verification", err)
	}

	return nil
}

func validateCredentials(email, password string) error {
	if !validation.ValidatePassword(password) {
		return status.Errorf(codes.InvalidArgument, "invalid password format")
	}

	if !validation.ValidateEmail(email) {
		return status.Errorf(codes.InvalidArgument, "invalid email format")
	}

	return nil
}

func (s *Service) CreateUserWithRole(ctx context.Context, email, password string, roleName role.RoleName) (*User, error) {
	var user *User
	err := pgTx.WithTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		hashedPassword, err := hasher.GenerateHashPassword(password)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}

		user, err = s.repo.createUser(ctx, tx, email, hashedPassword)
		if err != nil {
			return pgError.MapPostgresError("failed to create user", err)
		}

		err = s.roleService.SetUserRole(ctx, tx, user.ID, roleName)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*UserWithRoles, error) {
	err := validateCredentials(email, password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.getUser(ctx, email)
	if err != nil {
		return nil, pgError.MapPostgresError("failed to get user", err)
	}

	isCheck, err := hasher.CompareHashAndPassword(password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password check failed: %v", err)
	}

	if !isCheck {
		return nil, status.Error(codes.Unauthenticated, "password is incorrect")
	}

	if !user.IsVerified {
		return nil, status.Errorf(codes.PermissionDenied, "user is not verified")
	}

	roles, err := s.roleService.GetUserRolesName(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &UserWithRoles{
		UserID: user.ID,
		Roles:  roles,
	}, nil
}

func (s *Service) MarkUserIsVerified(ctx context.Context, tx pgx.Tx, userID int64) error {
	err := s.repo.markUserIsVerified(ctx, tx, userID)
	if err != nil {
		return pgError.MapPostgresError("failed to verify user", err)
	}

	return nil
}
