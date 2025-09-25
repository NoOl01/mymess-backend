package grpc_server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"profile/internal/grpc_client"
	"profile/internal/jwt"
	"proto/databasepb"
	pb "proto/profilepb"
	"results/succ"
	"strconv"
)

type Server struct {
	pb.UnimplementedProfileServiceServer
	Client databasepb.DatabaseServiceClient
}

func (s *Server) Update(_ context.Context, req *pb.UpdateRequest) (*pb.BaseResultResponse, error) {
	token := req.GetToken()

	userIdStr, err := jwt.ValidateJwt(token)
	if err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	resp, err := grpc_client.UpdateProfile(s.Client, req.GetValue(), req.GetType(), userId)
	if err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{Result: resp.Result}, nil
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
