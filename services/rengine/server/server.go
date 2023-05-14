package server

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/rengine/recommendation"
	"github.com/Inoi-K/Find-Me/services/rengine/selection"
	"github.com/Inoi-K/Find-Me/services/rengine/session"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedREngineServer
}

// Next gets the next user recommendation and returns it
func (s *server) Next(ctx context.Context, in *pb.NextRequest) (*pb.NextReply, error) {
	// create recommendations for the user if they do not exist yet
	if _, ok := session.SUS[in.SphereID][in.UserID]; !ok {
		usdt, err := database.GetUsersTag(ctx)
		if err != nil {
			log.Fatalf("failed to get user sphere tags %v", err)
		}
		matches, err := database.GetMatches(ctx, in.SphereID)
		if err != nil {
			log.Fatalf("failed to get matches %v", err)
		}
		w, err := database.GetWeights(ctx)
		if err != nil {
			log.Fatalf("failed to get weights %v", err)
		}
		searchFamiliar, err := database.GetSearchFamiliar(ctx, in.UserID, in.SphereID)
		if err != nil {
			log.Fatalf("failed to get search option %v", err)
		}

		// create recommendations for current user
		session.SUS[in.SphereID][in.UserID] = recommendation.CalculateSimilarities(in.UserID, in.SphereID, searchFamiliar, usdt, matches, w)
	}

	// no more recommendations
	if len(session.SUS[in.SphereID][in.UserID]) == 0 {
		return &pb.NextReply{
			NextUserID: -1,
		}, nil
	}

	// get the next recommendation
	nextID := selection.PickStrategy.Pick(session.SUS[in.SphereID][in.UserID])

	return &pb.NextReply{
		NextUserID: nextID,
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
