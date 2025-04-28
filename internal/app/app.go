package app

import (
	"context"
	"github.com/handmade-jewellery/user-service/internal/app/user"
	"github.com/handmade-jewellery/user-service/internal/server"
)

type App struct {
	server *server.Server
	impl   *user.Service
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	a.initDeps()

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	//ctx - нужен для запуска backgrounds processes
	return a.server.Run()
}

func (a *App) initDeps() {
	a.initImpl()
	a.initServer()
}

func (a *App) initImpl() {
	a.impl = user.NewService()
}

func (a *App) initServer() {
	a.server = server.NewServer(a.impl)
}
