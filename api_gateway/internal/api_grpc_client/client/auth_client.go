package client

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
		Nickname: username,
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

func RefreshRequest(client authpb.AuthServiceClient, accessToken string) (*authpb.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Refresh(ctx, &authpb.RefreshRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SendOtp(client authpb.AuthServiceClient, email string) (*authpb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SendOtp(ctx, &authpb.SendOtpRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ResetPassword(client authpb.AuthServiceClient, email string, code int32) (*authpb.ResetPasswordResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ResetPassword(ctx, &authpb.ResetPasswordRequest{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func UpdatePassword(client authpb.AuthServiceClient, email, password, resetToken string) (*authpb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.UpdatePassword(ctx, &authpb.UpdatePasswordRequest{
		Email:      email,
		Password:   password,
		ResetToken: resetToken,
	})
}

func MyProfile(client authpb.AuthServiceClient, token string) (*authpb.MyProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.MyProfile(ctx, &authpb.MyProfileRequest{
		Token: token,
	})
}
