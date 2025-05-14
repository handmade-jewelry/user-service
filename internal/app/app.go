package app

import (
	"github.com/handmade-jewelry/user-service/internal/app/user"
	"github.com/handmade-jewelry/user-service/internal/config"
	"github.com/handmade-jewelry/user-service/internal/server"
	"github.com/spf13/viper"
	"log"
	"time"
)

type App struct {
	cfg    *config.Config
	server *server.Server
	impl   *user.Service
}

func NewApp() (*App, error) {
	a := &App{}
	err := a.initDeps()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) initDeps() error {
	err := a.initConfig()
	if err != nil {
		return err
	}

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
		log.Fatalf("Failed to parse http gracefull timeout duration config: %v", err)
		return err
	}

	a.cfg = &config.Config{
		GRPCPort:            viper.GetString(config.GRPCPort),
		GRPCNetwork:         viper.GetString(config.GRPCNetwork),
		HTTPPort:            viper.GetString(config.HTTPPort),
		HTTPHost:            viper.GetString(config.HTTPHost),
		HTTPGracefulTimeout: httpGracefulTimeout,
	}

	return nil
}

func (a *App) initImpl() {
	a.impl = user.NewService()
}

func (a *App) initServer() {
	opts := &server.Opts{
		GrpcPort:        a.cfg.GRPCPort,
		GrpcNetwork:     a.cfg.GRPCNetwork,
		HttpPort:        a.cfg.HTTPPort,
		HttpHost:        a.cfg.HTTPHost,
		GracefulTimeout: a.cfg.HTTPGracefulTimeout,
	}
	a.server = server.NewServer(a.impl, opts)
}
