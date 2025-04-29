package main

import (
	"github.com/handmade-jewellery/user-service/internal/app"
	"log"
)

func main() {
	//ctx := context.Background() todo
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
