package client

import (
	"context"
	"proto/databasepb"
	"time"
)

func FindUser(client databasepb.DatabaseServiceClient, username string) (*databasepb.FindProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.FindProfile(ctx, &databasepb.FindProfileRequest{
		Name: username,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetProfileInfo(client databasepb.DatabaseServiceClient, id int64) (*databasepb.GetProfileInfoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetProfileInfo(ctx, &databasepb.GetProfileInfoRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
