package client

import (
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Match     pb.MatchClient
	connMatch *grpc.ClientConn
)

// Open creates connections to other services
func Open() error {
	var err error

	// Set up a connection to the match server
	address := config.C.MatchHost + ":" + config.C.MatchPort
	connMatch, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	Match = pb.NewMatchClient(connMatch)

	return nil
}

// Close closes connections to other services
func Close() {
	connMatch.Close()
}
