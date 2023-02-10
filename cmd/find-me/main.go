package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/internal/database"
	"github.com/Inoi-K/Find-Me/internal/user"
	"log"
)

func main() {
	ctx := context.Background()
	_, err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db")
	}

	users, err := database.GetUsers(ctx)
	if err != nil {
		log.Fatalf("failed to get users %v", err)
	}

	user.ShowSimilarityAll(users, "work")
}
