package grpc_client

import (
	"auth/internal/storage"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	databasebp "proto/databasepb"
	"results/errs"
)

func GrpcClientConnect() (databasebp.DatabaseServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", storage.Env.DatabaseHost, storage.Env.DatabasePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%s: %v", errs.GrpcClientConnectFailed, err)
	}

	return databasebp.NewDatabaseServiceClient(gRpcConn), gRpcConn, nil
}
