package api_grpc_client

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"proto/authpb"
	"proto/databasepb"
	"proto/smtppb"
	"time"
)

func AuthPing(client authpb.AuthServiceClient) (*authpb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DatabasePing(client databasepb.DatabaseServiceClient) (*databasepb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SmtpPing(client smtppb.SmtpServiceClient) (*smtppb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
