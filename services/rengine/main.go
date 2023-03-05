package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/rengine/server"
	"log"
)

func main() {
	config.ReadConfig()

	ctx := context.Background()

	err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	server.Start()
}
