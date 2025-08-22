package grpc_client

import (
	"chat_kafka_broker/internal/storage"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"proto/scyllapb"
	"results/errs"
)

func GrpcScyllaClientConnect() (scyllapb.ScyllaServiceClient, *grpc.ClientConn, error) {
	gRpcConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", storage.Env.DatabaseHost, storage.Env.DatabasePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%s: %v", errs.GrpcClientConnectFailed, err)
	}

	return scyllapb.NewScyllaServiceClient(gRpcConn), gRpcConn, nil
}
