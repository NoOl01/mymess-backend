package grpc_server

import (
	"auth/internal/grpc_client"
	"auth/internal/jwt"
	"context"
	"fmt"
	pb "proto/authpb"
	"proto/databasepb"
	"results/errs"
	"results/succ"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	Client databasepb.DatabaseServiceClient
}

func (s *Server) Register(_ context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	username := req.GetUsername()
	email := req.GetEmail()
	password := req.GetPassword()

	if username == "" || email == "" || password == "" {
		return &pb.AuthResponse{
			AccessToken: "",
			Error:       "Username, Email and Password is empty",
		}, nil
	}

	resp, err := grpc_client.CreateNewUser(s.Client, username, email, password)
	if err != nil {
		return nil, err
	}

	if resp.Result != succ.RecordCreated {
		return nil, fmt.Errorf(resp.Result)
	}

	token, err := jwt.GenerateToken(resp.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		AccessToken: token,
		Error:       "",
	}, nil
}

func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	var dbUsername, dbEmail string

	switch login := req.LoginMethod.(type) {
	case *pb.LoginRequest_Email:
		dbEmail = login.Email
	case *pb.LoginRequest_Username:
		dbUsername = login.Username
	default:
		return &pb.AuthResponse{Error: errs.InvalidRequestBody.Error()}, errs.InvalidRequestBody
	}

	resp, err := grpc_client.CheckUser(s.Client, dbUsername, dbEmail, req.GetPassword())
	if err != nil {
		return &pb.AuthResponse{Error: err.Error()}, err
	}

	if resp.Result != succ.Ok {
		return &pb.AuthResponse{Error: resp.Result}, fmt.Errorf(resp.Result)
	}

	token, err := jwt.GenerateToken(resp.UserId)
	if err != nil {
		return &pb.AuthResponse{Error: err.Error()}, err
	}

	return &pb.AuthResponse{AccessToken: token}, nil
}
