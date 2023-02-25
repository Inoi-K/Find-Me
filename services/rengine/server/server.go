package server

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/rengine/recommendation"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedREngineServer
}

func (s *server) UpdateRecommendations(ctx context.Context, in *pb.UpdateRecommendationsRequest) (*pb.UpdateRecommendationsReply, error) {
	ust, err := database.GetUserSphereTag(ctx)
	if err != nil {
		log.Fatalf("failed to get user sphere tags %v", err)
	}

	// TODO create recommendations for current user
	recommendations := recommendation.CreateRecommendationsForUser(in.UserID, ust)

	// TODO update recommendations of affected online users

	//recommendation.ShowSimilarityAll(users, "work")

	return &pb.UpdateRecommendationsReply{
		RecommendationIDs: recommendations,
	}, nil
}

func Start() {
	lis, err := net.Listen("tcp", ":"+config.C.REnginePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterREngineServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
