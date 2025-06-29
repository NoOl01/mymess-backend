package api_grpc_client

import (
	"api_gateway/internal/api_storage"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "proto/authpb"
	"results/errs"
)

func GrpcClientConnect() (pb.AuthServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("dns:///%s:%s", api_storage.Env.AuthHost, api_storage.Env.AuthPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", errs.GrpcClientConnectFailed, err)
	}

	return pb.NewAuthServiceClient(gRpcConn), gRpcConn, nil
}
