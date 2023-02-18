package client

import (
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Profile pb.ProfileClient
	conn    *grpc.ClientConn
)

// Open creates connections to other services
func Open() error {
	var err error
	// Set up a connection to the server.
	address := config.C.ProfileHost + ":" + config.C.ProfilePort
	conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	//defer conn.Close()
	Profile = pb.NewProfileClient(conn)
	return nil
}

// Close closes connections to other services
func Close() {
	conn.Close()
}
