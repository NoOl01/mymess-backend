package routes

import (
	"context"
	"fmt"
	pb "proto/databasepb"
)

type Server struct {
	pb.UnimplementedDatabaseServiceServer
}

func (s *Server) Register(_ context.Context, req *pb.CreateNewUserRequest) (*pb.CreateNewUserResponse, error) {
	username := req.GetUsername()
	email := req.GetEmail()
	password := req.GetPassword()

	if username == "" || email == "" || password == "" {
		return &pb.CreateNewUserResponse{
			Result: "Missing fields",
		}, nil
	}

	return &pb.CreateNewUserResponse{
		Result: "Ok",
	}, nil
}

func (s *Server) LoginUser(_ context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	var userNameOrEmail string

	switch login := req.LoginMethod.(type) {
	case *pb.LoginUserRequest_Username:
		userNameOrEmail = login.Username
	case *pb.LoginUserRequest_Email:
		userNameOrEmail = login.Email
	default:
		return &pb.LoginUserResponse{
			Result: "Username or Email required",
		}, nil
	}

	return &pb.LoginUserResponse{
		Result: fmt.Sprintf("User: %s", userNameOrEmail),
	}, nil
}
