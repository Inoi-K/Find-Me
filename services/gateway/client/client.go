package client

import (
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Profile pb.ProfileClient
	conn    *grpc.ClientConn
)

func Open() error {
	var err error
	// Set up a connection to the server.
	address := fmt.Sprintf("%s:%s", config.C.ProfileHost, config.C.ProfilePort)
	conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	//defer conn.Close()
	Profile = pb.NewProfileClient(conn)
	return nil
}

func Close() {
	conn.Close()
}
