package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/rengine/server"
	"github.com/Inoi-K/Find-Me/services/rengine/session"
	"log"
)

func main() {
	config.ReadConfig()

	ctx := context.Background()

	err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	session.SUR = make(map[int64]map[int64][]int64)
	session.SUR[config.C.SphereID] = make(map[int64][]int64)

	server.Start()
}
