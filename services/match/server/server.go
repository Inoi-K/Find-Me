package server

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/database"
	"github.com/Inoi-K/Find-Me/services/match/client"
	"github.com/Inoi-K/Find-Me/services/match/session"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMatchServer
}

// Next gets the next user recommendation and returns it
func (s *server) Next(ctx context.Context, in *pb.NextRequest) (*pb.NextReply, error) {
	log.Printf("Received next %v", in.UserID)

	// create recommendations for the user if they do not exist yet
	if _, ok := session.SUR[in.SphereID][in.UserID]; !ok {
		err := updateUserRecommendations(ctx, in.UserID, in.SphereID)
		if err != nil {
			return nil, err
		}
	}

	// no more recommendations
	if len(session.SUR[in.SphereID][in.UserID]) == 0 {
		return &pb.NextReply{
			NextUserID: -1,
		}, nil
	}

	// get the next recommendation
	nextID := session.SUR[in.SphereID][in.UserID][0]

	return &pb.NextReply{
		NextUserID: nextID,
	}, nil
}

// UpdateRecommendations updates the recommendations of the user and affected ones
func (s *server) UpdateRecommendations(ctx context.Context, in *pb.UpdateRecommendationsRequest) (*pb.MatchEmpty, error) {
	log.Printf("Received update recommendations: %v", in.UserID)

	err := updateUserRecommendations(ctx, in.UserID, in.SphereID)
	if err != nil {
		return nil, err
	}

	// TODO update recommendations of affected online users

	return &pb.MatchEmpty{}, nil
}

// updateUserRecommendations updates recommendations of the user per sphere
func updateUserRecommendations(ctx context.Context, userID, sphereID int64) error {
	var err error
	session.SUR[sphereID][userID], err = getRecommendations(ctx, userID, sphereID)
	if err != nil {
		return err
	}
	return nil
}

// getRecommendations get the recommendations of user
func getRecommendations(ctx context.Context, userID, sphereID int64) ([]int64, error) {
	// contact the rengine service to get recommendations
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

// Match handles user's like reaction
func (s *server) Match(ctx context.Context, in *pb.MatchRequest) (*pb.MatchReply, error) {
	log.Printf("Received match: %v", in.FromID)

	// remove the recommendation from the slice
	if in.ToID == session.SUR[in.SphereID][in.FromID][0] {
		session.SUR[in.SphereID][in.FromID] = session.SUR[in.SphereID][in.FromID][1:]
	} else {
		log.Printf("strange behaviour: received like/dislike fromID=%d toID=%d in sphereID=%d, but the %d user recommendation was give in the last time", in.FromID, in.ToID, in.SphereID, session.SUR[in.SphereID][in.FromID][0])
	}

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
