package api_grpc_client

import (
	"api_gateway/internal/api_storage"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "proto/authpb"
	"proto/databasepb"
	"proto/smtppb"
	"results/errs"
)

// todo Этот ужас тоже переделать

func GrpcAuthClientConnect() (pb.AuthServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("dns:///%s:%s", api_storage.Env.AuthHost, api_storage.Env.AuthPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", errs.GrpcClientConnectFailed, err)
	}

	return pb.NewAuthServiceClient(gRpcConn), gRpcConn, nil
}

func GrpcDatabaseClientConnect() (databasepb.DatabaseServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("dns:///%s:%s", api_storage.Env.DbHost, api_storage.Env.DbPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", errs.GrpcClientConnectFailed, err)
	}

	return databasepb.NewDatabaseServiceClient(gRpcConn), gRpcConn, nil
}

func GrpcSmtpClientConnect() (smtppb.SmtpServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("dns:///%s:%s", api_storage.Env.SmtpServiceHost, api_storage.Env.SmtpServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", errs.GrpcClientConnectFailed, err)
	}

	return smtppb.NewSmtpServiceClient(gRpcConn), gRpcConn, nil
}
