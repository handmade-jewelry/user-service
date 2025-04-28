package main

import (
	"context"
	"github.com/handmade-jewellery/user-service/internal/app"
	"log"
)

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
