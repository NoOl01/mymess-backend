package api_grpc_client

import (
	"api_gateway/internal/api_storage"
	"errs"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "proto/authpb"
)

func GrpcClientConnect() (pb.AuthServiceClient, error) {
	jRpcConn, err := grpc.NewClient(
		fmt.Sprintf("dns:///localhost:%s", api_storage.Env.AuthPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.GrpcClientConnectFailed, err)
	}

	defer func(jRpcConn *grpc.ClientConn) {
		err := jRpcConn.Close()
		if err != nil {
			fmt.Printf("%s, %v \n", errs.GrpcClientCloseFailed, err)
			return
		}
	}(jRpcConn)

	return pb.NewAuthServiceClient(jRpcConn), nil
}
