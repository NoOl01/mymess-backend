package api_grpc_client

import (
	"context"
	"proto/authpb"
	"results/errs"
	"time"
)

func RegisterRequest(client authpb.AuthServiceClient, username, email, password string) (*authpb.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &authpb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func LoginRequest(client authpb.AuthServiceClient, username, email, password string) (*authpb.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loginReq *authpb.LoginRequest

	switch {
	case email != "":
		loginReq = &authpb.LoginRequest{
			LoginMethod: &authpb.LoginRequest_Email{Email: email},
			Password:    password,
		}
	case username != "":
		loginReq = &authpb.LoginRequest{
			LoginMethod: &authpb.LoginRequest_Username{Username: username},
			Password:    password,
		}
	default:
		return nil, errs.InvalidRequestBody
	}

	resp, err := client.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
