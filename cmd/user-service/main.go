package main

import (
	"context"
	"log"

	"github.com/handmade-jewelry/user-service/internal/app"
	"github.com/handmade-jewelry/user-service/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to create app: %v", err)
	}

	err = a.Run()
	if err != nil {
		logger.Fatal("failed to run app: %v", err)
	}
}
