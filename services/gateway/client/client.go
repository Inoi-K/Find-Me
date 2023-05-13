package client

import (
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Profile     pb.ProfileClient
	Match       pb.MatchClient
	REngine     pb.REngineClient
	connProfile *grpc.ClientConn
	connMatch   *grpc.ClientConn
	connREngine *grpc.ClientConn
)

// Open creates connections to other services
func Open() error {
	var err error

	// Set up a connection to the profile server.
	address := config.C.ProfileHost + ":" + config.C.ProfilePort
	connProfile, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	Profile = pb.NewProfileClient(connProfile)

	// Set up a connection to the match server
	address = config.C.MatchHost + ":" + config.C.MatchPort
	connMatch, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	Match = pb.NewMatchClient(connMatch)

	// Set up a connection to the rengine server
	address = config.C.REngineHost + ":" + config.C.REnginePort
	connREngine, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	REngine = pb.NewREngineClient(connREngine)

	return nil
}

// Close closes connections to other services
func Close() {
	connProfile.Close()
	connMatch.Close()
}
