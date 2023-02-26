package server

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/match/client"
	"github.com/Inoi-K/Find-Me/services/match/session"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMatchServer
}

func (s *server) Next(ctx context.Context, in *pb.NextRequest) (*pb.NextReply, error) {
	var nextID int64
	if _, ok := session.SUR[in.SphereID][in.UserID]; !ok {
		// get recommendations for current user
		var err error
		session.SUR[in.SphereID][in.UserID], err = getRecommendations(ctx, in.UserID, in.SphereID)
		if err != nil {
			return nil, err
		}
	}

	nextID = session.SUR[in.SphereID][in.UserID][0]
	session.SUR[in.SphereID][in.UserID] = session.SUR[in.SphereID][in.UserID][1:]

	return &pb.NextReply{
		NextUserID: nextID,
	}, nil
}

func (s *server) UpdateRecommendations(ctx context.Context, in *pb.UpdateRecommendationsRequest) (*pb.MatchEmpty, error) {
	// update recommendations of current user
	var err error
	session.SUR[in.SphereID][in.UserID], err = getRecommendations(ctx, in.UserID, in.SphereID)
	if err != nil {
		return nil, err
	}

	// TODO update recommendations of affected online users

	return &pb.MatchEmpty{}, nil
}

func getRecommendations(ctx context.Context, userID, sphereID int64) ([]int64, error) {
	// contact the rengine server
	ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
	defer cancel()
	rep, err := client.REngine.GetRecommendations(ctx2, &pb.GetRecommendationsRequest{
		UserID:   userID,
		SphereID: sphereID,
	})
	if err != nil {
		return nil, err
	}
	return rep.RecommendationIDs, nil
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
