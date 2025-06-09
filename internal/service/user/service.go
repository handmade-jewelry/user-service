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
	pgUtils "github.com/handmade-jewelry/user-service/libs/pgutils"
	"github.com/handmade-jewelry/user-service/logger"
)

type Service struct {
	repo                *repository
	dbPool              *pgxpool.Pool
	roleService         *role.Service
	verificationService *verification.Service
}

func NewService(dbPool *pgxpool.Pool, roleService *role.Service, verificationService *verification.Service) *Service {
	return &Service{
		repo:                newRepository(dbPool),
		dbPool:              dbPool,
		roleService:         roleService,
		verificationService: verificationService,
	}
}

func (s *Service) RegisterSeller(ctx context.Context, email, password string) error {
	return s.Register(ctx, email, password, role.SellerRoleName)
}

func (s *Service) RegisterCustomer(ctx context.Context, email, password string) error {
	return s.Register(ctx, email, password, role.CustomerRoleName)
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
		logger.ErrorWithFields("failed to send verification", err, "user_id", user.ID)
	}

	return nil
}

func validateCredentials(email, password string) error {
	if !validation.ValidatePassword(password) {
		return status.Error(codes.InvalidArgument, "invalid password format")
	}

	if !validation.ValidateEmail(email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	return nil
}

func (s *Service) CreateUserWithRole(ctx context.Context, email, password string, roleName role.RoleName) (*User, error) {
	var user *User
	err := pgUtils.WithTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		hashedPassword, err := hasher.GenerateHashPassword(password)
		if err != nil {
			logger.Error("failed to hash password", err)
			return status.Errorf(codes.Internal, "internal error")
		}

		user, err = s.repo.createUser(ctx, tx, email, hashedPassword)
		if err != nil {
			return pgUtils.MapPostgresError("user", err)
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
		logger.ErrorWithFields("failed to get user", pgUtils.MapPostgresError("", err), "email", email)
		return nil, pgUtils.MapPostgresError("user", err)
	}

	isCheck, err := hasher.CompareHashAndPassword(password, user.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "email or password is incorrect")
	}

	if !isCheck {
		return nil, status.Error(codes.Unauthenticated, "email or password is incorrect")
	}

	if !user.IsVerified {
		return nil, status.Errorf(codes.PermissionDenied, "user is not verified")
	}

	roles, err := s.roleService.GetUserRolesName(ctx, user.ID)
	if err != nil {
		logger.ErrorWithFields("failed to fetch roles", err, "user_id", user.ID)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &UserWithRoles{
		UserID: user.ID,
		Roles:  roles,
	}, nil
}

func (s *Service) MarkUserIsVerified(ctx context.Context, tx pgx.Tx, userID int64) error {
	err := s.repo.markUserIsVerified(ctx, tx, userID)
	if err != nil {
		return pgUtils.MapPostgresError("user", err)
	}

	return nil
}
