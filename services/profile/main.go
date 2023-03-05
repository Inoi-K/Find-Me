package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/profile/client"
	"github.com/Inoi-K/Find-Me/services/profile/server"
	"log"
)

func main() {
	config.ReadConfig()

	ctx := context.Background()

	err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	err = client.Open()
	if err != nil {
		log.Fatalf("couldn't connect to grpc client: %v", err)
	}

	server.Start()
}
