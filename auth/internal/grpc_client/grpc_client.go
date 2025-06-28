package grpc_client

import (
	"context"
	"proto/databasepb"
	"time"
)

func CreateNewUser(client databasepb.DatabaseServiceClient, username, email, password string) (*databasepb.CreateNewUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &databasepb.CreateNewUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func CheckUser(client databasepb.DatabaseServiceClient, username, email, password string) (*databasepb.CreateNewUserResponse, error) {
	// todo
	return nil, nil
}
