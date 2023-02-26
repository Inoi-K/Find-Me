package client

import (
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	REngine     pb.REngineClient
	connREngine *grpc.ClientConn
)

// Open creates connections to other services
func Open() error {
	var err error

	// Set up a connection to the rengine server
	address := config.C.REngineHost + ":" + config.C.REnginePort
	connREngine, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	REngine = pb.NewREngineClient(connREngine)

	return nil
}

// Close closes connections to other services
func Close() {
	connREngine.Close()
}
