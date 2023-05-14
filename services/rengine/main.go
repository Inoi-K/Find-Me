package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/rengine/selection"
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

	session.SUS = make(map[int64]map[int64]map[int64]float64)
	session.SUS[config.C.SphereID] = make(map[int64]map[int64]float64)

	selection.PickStrategy = selection.TournamentStrategy{PopulationCount: 3}

	server.Start()
}
