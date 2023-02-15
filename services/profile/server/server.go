package server

import (
	"context"
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedProfileServer
}

func (s *server) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpReply, error) {
	log.Printf("Received signup: %v", in.GetName())

	err := database.AddUser(ctx, in)
	if err != nil {
		return &pb.SignUpReply{IsOk: false}, err
	}

	return &pb.SignUpReply{IsOk: true}, nil
}

func (s *server) Exists(ctx context.Context, in *pb.ExistsRequest) (*pb.ExistsReply, error) {
	log.Printf("Received exists: %v", in.GetUserID())

	exists, err := database.UserExists(ctx, in.GetUserID())
	if err != nil {
		return nil, err
	}
	log.Println(exists)
	return &pb.ExistsReply{Exists: exists}, nil
}

func Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.C.ProfilePort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProfileServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
