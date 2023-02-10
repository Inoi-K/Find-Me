package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/recommendations/recommendation"
	"log"
)

func main() {
	config.ReadConfig()

	ctx := context.Background()
	_, err := database.ConnectDB(ctx, config.C.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	users, err := database.GetUsers(ctx)
	if err != nil {
		log.Fatalf("failed to get users %v", err)
	}

	recommendation.ShowSimilarityAll(users, "work")
}
