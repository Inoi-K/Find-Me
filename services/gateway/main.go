package main

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/client"
	"github.com/Inoi-K/Find-Me/services/gateway/handler"
	"log"
)

func main() {
	config.ReadConfig()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	err := client.Open()
	if err != nil {
		log.Fatalf("culdn't connect to grpc client: %v", err)
	}

	err = handler.Start(ctx)
	if err != nil {
		log.Fatalf("couldn't start handler %v", err)
	}

	// Tell the user the bot is online
	log.Println("Start listening for updates...")

	select {}

	// Wait for a newline symbol, then cancel handling updates
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}
