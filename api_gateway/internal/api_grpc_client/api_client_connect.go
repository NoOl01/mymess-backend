package api_grpc_client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"results/errs"
)

type ClientFactory[T any] func(conn grpc.ClientConnInterface) T

type ClientParams[T any] struct {
	ServiceHost string
	ServicePort string
	Factory     ClientFactory[T]
}

func GrpcClientConnect[T any](params ClientParams[T]) (T, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("dns:///%s:%s", params.ServiceHost, params.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		var zero T
		return zero, nil, fmt.Errorf("%w: %s:%s. %v\n", errs.GrpcClientConnectFailed, params.ServiceHost, params.ServicePort, err)
	}

	client := params.Factory(conn)
	return client, conn, nil
}

func Connect[T any](host, port string, factory func(grpc.ClientConnInterface) T) ClientParams[T] {
	return ClientParams[T]{
		ServiceHost: host,
		ServicePort: port,
		Factory:     factory,
	}
}
