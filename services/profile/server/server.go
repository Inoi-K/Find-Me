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
	pb.UnimplementedProfileServer
}

// SignUp adds user to database
func (s *server) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.ProfileEmpty, error) {
	log.Printf("Received sign up: %v", in.UserID)

	err := database.AddUser(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.ProfileEmpty{}, nil
}

// GetUserMain gets main info of a user
func (s *server) GetUserMain(ctx context.Context, in *pb.GetUserMainRequest) (*pb.GetUserMainReply, error) {
	log.Printf("Received get user: %v", in.UserID)

	user, err := database.GetUserMain(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserMainReply{
		Name:       user.Name,
		Gender:     user.Gender,
		Age:        user.Age,
		Faculty:    user.Faculty,
		University: user.University,
		Username:   user.Username,
	}, nil
}

// GetUserAdditional gets an additional info of a user
func (s *server) GetUserAdditional(ctx context.Context, in *pb.GetUserAdditionalRequest) (*pb.GetUserAdditionalReply, error) {
	log.Printf("Received get user sphere: %v", in.UserID)

	sphereInfo, err := database.GetUserAdditional(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	// convert tags map to slice
	tags := make([]string, len(sphereInfo[in.SphereID].Tags))
	i := 0
	for tag := range sphereInfo[in.SphereID].Tags {
		tags[i] = tag
		i++
	}

	return &pb.GetUserAdditionalReply{
		Description: sphereInfo[in.SphereID].Description,
		PhotoID:     sphereInfo[in.SphereID].PhotoID,
		Tags:        tags,
	}, nil
}

// Exists checks if user exists in database
func (s *server) Exists(ctx context.Context, in *pb.ExistsRequest) (*pb.ExistsReply, error) {
	log.Printf("Received exists: %v", in.UserID)

	exists, err := database.UserExists(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.ExistsReply{Exists: exists}, nil
}

// Edit updates field of a user with a new value
func (s *server) Edit(ctx context.Context, in *pb.EditRequest) (*pb.ProfileEmpty, error) {
	log.Printf("Received edit field: %v", in.UserID)

	switch in.Field {
	case AgeField, FacultyField, PhotoField, DescriptionField:
		switch n := len(in.Value); {
		case n == 0:
			return nil, WrongArgumentsNumberError
		case n > 1:
			log.Printf("strange behaviour in EditField: required 1 value, but %d given - %v", n, in.Value)
		}

		err := database.EditField(ctx, in.Field, in.Value[0], in.UserID, in.SphereID)
		if err != nil {
			return nil, err
		}

	case TagsField:
		err := database.EditTags(ctx, in.Value, in.UserID, in.SphereID)
		if err != nil {
			return nil, err
		}

	default:
		return nil, UnknownFieldError
	}

	// contact the match server to update recommendations
	//ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
	//defer cancel()
	//_, err := client.Match.UpdateRecommendations(ctx2, &pb.UpdateRecommendationsRequest{
	//	UserID:   in.UserID,
	//	SphereID: in.SphereID,
	//})
	//if err != nil {
	//	return nil, err
	//}

	return &pb.ProfileEmpty{}, nil
}

// GetTags returns tags of the sphere
func (s *server) GetTags(ctx context.Context, in *pb.GetTagsRequest) (*pb.GetTagsReply, error) {
	log.Printf("Received get tags: %v", in.SphereID)

	// get tags
	tags, err := database.GetTags(ctx, in.SphereID)
	if err != nil {
		return nil, err
	}

	// convert tag model to 2 slices
	rep := &pb.GetTagsReply{
		TagNames: make([]string, len(tags)),
		TagIDs:   make([]string, len(tags)),
	}
	for i := 0; i < len(tags); i++ {
		rep.TagNames[i] = tags[i].Name
		rep.TagIDs[i] = tags[i].ID
	}

	return rep, nil
}

func Start() {
	lis, err := net.Listen("tcp", ":"+config.C.ProfilePort)
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
