package grpc_client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"profile/internal/storage"
	databasebp "proto/databasepb"
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
