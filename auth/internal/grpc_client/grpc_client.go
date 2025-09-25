package grpc_client

import (
	"context"
	"proto/databasepb"
	"proto/smtppb"
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

func SendOtp(client smtppb.SmtpServiceClient, email string, code int32) (*smtppb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.SendOtp(ctx, &smtppb.SendOtpRequest{
		Email: email,
		Code:  code,
	})
}

func UpdatePassword(client databasepb.DatabaseServiceClient, email, password string) (*databasepb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.UpdatePassword(ctx, &databasepb.UpdatePasswordRequest{
		Email:    email,
		Password: password,
	})
}

func MyProfile(client databasepb.DatabaseServiceClient, id int64) (*databasepb.MyProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.MyProfile(ctx, &databasepb.GetProfileInfoRequest{
		Id: id,
	})
}
