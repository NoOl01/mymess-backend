package grpc_client

import (
	"context"
	"proto/databasepb"
	"time"
)

func UpdateProfile(client databasepb.DatabaseServiceClient, value, updateType string, userId int64) (*databasepb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.UpdateProfile(ctx, &databasepb.UpdateRequest{
		Value:  value,
		Type:   updateType,
		UserId: userId,
	})
}
