package grpc_client

import (
	"auth/internal/storage"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	databasebp "proto/databasepb"
	"proto/smtppb"
	"results/errs"
)

func GrpcDatabaseClientConnect() (databasebp.DatabaseServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", storage.Env.DatabaseHost, storage.Env.DatabasePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%s: %v", errs.GrpcClientConnectFailed, err)
	}

	return databasebp.NewDatabaseServiceClient(gRpcConn), gRpcConn, nil
}

func GrpcSmtpClientConnect() (smtppb.SmtpServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", storage.Env.SmtpHost, storage.Env.SmtpPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%s: %v", errs.GrpcClientConnectFailed, err)
	}

	return smtppb.NewSmtpServiceClient(gRpcConn), gRpcConn, nil
}
