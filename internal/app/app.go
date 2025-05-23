package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	userServer "github.com/handmade-jewelry/user-service/internal/app/user"
	"github.com/handmade-jewelry/user-service/internal/config"
	"github.com/handmade-jewelry/user-service/internal/server"
	"github.com/handmade-jewelry/user-service/internal/service/role"
	"github.com/handmade-jewelry/user-service/internal/service/user"
	userVerification "github.com/handmade-jewelry/user-service/internal/service/user-verification"
	"github.com/handmade-jewelry/user-service/internal/service/verification"
	"github.com/handmade-jewelry/user-service/logger"
)

type App struct {
	cfg                     *config.Config
	userServiceServer       *userServer.UserServiceServer
	server                  *server.Server
	dBPool                  *pgxpool.Pool
	userService             *user.Service
	roleService             *role.Service
	verificationService     *verification.Service
	userVerificationService *userVerification.Service
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) initDeps(ctx context.Context) error {
	err := a.initConfig()
	if err != nil {
		return err
	}

	err = a.initDb(ctx)
	if err != nil {
		return err
	}

	a.initService()
	a.initImpl()
	a.initServer()

	return nil
}

func (a *App) initConfig() error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	httpGracefulTimeout, err := time.ParseDuration(viper.GetString(config.HTTPGracefulTimeout))
	if err != nil {
		return fmt.Errorf("failed to parse http gracefull timeout duration config: %w", err)
	}

	dBMaxConLifetime, err := time.ParseDuration(viper.GetString(config.DBMaxConLifetime))
	if err != nil {
		return fmt.Errorf("failed to parse dBPool max conns lifetime duration config: %w", err)
	}

	verificationTokenExp, err := time.ParseDuration(viper.GetString(config.VerificationTokenExp))
	if err != nil {
		return fmt.Errorf("failed to parse dBPool varification token exp duration config: %w", err)
	}

	a.cfg = &config.Config{
		GRPCPort:             viper.GetString(config.GRPCPort),
		GRPCNetwork:          viper.GetString(config.GRPCNetwork),
		HTTPPort:             viper.GetString(config.HTTPPort),
		HTTPHost:             viper.GetString(config.HTTPHost),
		HTTPGracefulTimeout:  httpGracefulTimeout,
		DBName:               viper.GetString(config.DBName),
		DBUser:               viper.GetString(config.DBUser),
		DBPassword:           viper.GetString(config.DBPassword),
		DbHost:               viper.GetString(config.DBHost),
		DbPort:               viper.GetUint16(config.DBPort),
		SSLMode:              viper.GetString(config.SSLMode),
		DBMaxCons:            viper.GetInt32(config.DBMaxCons),
		DBMinCons:            viper.GetInt32(config.DBMinCons),
		DBMaxConLifetime:     dBMaxConLifetime,
		VerificationTokenExp: verificationTokenExp,
	}

	return nil
}

func (a *App) initImpl() {
	a.userServiceServer = userServer.NewUserServiceServer(
		a.userService,
		a.roleService,
		a.verificationService,
		a.userVerificationService,
	)
}

func (a *App) initServer() {
	opts := &server.Opts{
		GrpcPort:        a.cfg.GRPCPort,
		GrpcNetwork:     a.cfg.GRPCNetwork,
		HttpPort:        a.cfg.HTTPPort,
		HttpHost:        a.cfg.HTTPHost,
		GracefulTimeout: a.cfg.HTTPGracefulTimeout,
	}
	a.server = server.NewServer(a.userServiceServer, opts)
}

func (a *App) initDb(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		a.cfg.DBUser,
		a.cfg.DBPassword,
		a.cfg.DbHost,
		a.cfg.DbPort,
		a.cfg.DBName,
		a.cfg.SSLMode,
	)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse db config: %w", err)
	}

	cfg.MaxConns = a.cfg.DBMaxCons
	cfg.MinConns = a.cfg.DBMinCons
	cfg.MaxConnLifetime = a.cfg.DBMaxConLifetime

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("unable to create pool: %w", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	a.dBPool = dbPool

	logger.Info(
		"Database connection established",
		a.cfg.DbHost,
		strconv.Itoa(int(a.cfg.DbPort)))

	return nil
}

func (a *App) initService() {
	a.roleService = role.NewService(a.dBPool)
	a.verificationService = verification.NewService(a.dBPool, a.cfg.VerificationTokenExp)
	a.userService = user.NewService(a.dBPool, a.roleService, a.verificationService)
	a.userVerificationService = userVerification.NewService(a.userService, a.verificationService, a.dBPool)
}
