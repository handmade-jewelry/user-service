package app

import (
	"github.com/handmade-jewelry/user-service/internal/app/user"
	"github.com/handmade-jewelry/user-service/internal/server"
	"github.com/handmade-jewelry/user-service/util"
	"github.com/spf13/viper"
	"log"
	"time"
)

type App struct {
	server     *server.Server
	impl       *user.Service
	jwtService *util.JWTService
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
	err := initConfig()
	if err != nil {
		return err
	}

	a.initUtils()
	a.initImpl()
	a.initServer()

	return nil
}

func (a *App) initImpl() {
	a.impl = user.NewService()
}

func (a *App) initServer() {
	a.server = server.NewServer(a.impl)
}

func (a *App) initUtils() {
	authTokenExp, err := time.ParseDuration(viper.GetString(authTokenExpMin))
	if err != nil {
		log.Fatalf("Ошибка при парсинге длительности auth: %v", err)
	}

	refreshTokenExp, err := time.ParseDuration(viper.GetString(refreshTokenExpMin))
	if err != nil {
		log.Fatalf("Ошибка при парсинге длительности refresh: %v", err)
	}

	a.jwtService = util.NewJWTService(viper.GetString(jwtTokenSecret), authTokenExp, refreshTokenExp)
}
