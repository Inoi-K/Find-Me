package server

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMatchServer
}

// Match handles user's like reaction
func (s *server) Match(ctx context.Context, in *pb.MatchRequest) (*pb.MatchReply, error) {
	log.Printf("Received match: %v", in.FromID)

	isReciprocated, err := database.Match(ctx, in.FromID, in.ToID, in.SphereID, in.IsLike)
	if err != nil {
		return nil, err
	}

	return &pb.MatchReply{
		IsReciprocated: isReciprocated,
	}, nil
}

func Start() {
	lis, err := net.Listen("tcp", ":"+config.C.MatchPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMatchServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
