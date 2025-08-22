package client

import (
	"context"
	"proto/profilepb"
	"time"
)

func UpdateProfile(client profilepb.ProfileServiceClient, value, updateType, token string) (*profilepb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Update(ctx, &profilepb.UpdateRequest{
		Value: value,
		Type:  updateType,
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
