package grpc_client

import (
	"context"
	"proto/databasepb"
	"results/errs"
	"time"
)

func CreateNewUser(client databasepb.DatabaseServiceClient, nickname, email, password string) (*databasepb.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Register(ctx, &databasepb.CreateNewUserRequest{
		Nickname: nickname,
		Email:    email,
		Password: password,
	})
}

func CheckUser(client databasepb.DatabaseServiceClient, username, email, password string) (*databasepb.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &databasepb.LoginUserRequest{
		Password: password,
	}

	switch {
	case username != "":
		req.LoginMethod = &databasepb.LoginUserRequest_Username{Username: username}
	case email != "":
		req.LoginMethod = &databasepb.LoginUserRequest_Email{Email: email}
	default:
		return nil, errs.InvalidRequestBody
	}

	return client.Login(ctx, req)
}
