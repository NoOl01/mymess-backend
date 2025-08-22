package grpc_client

import (
	"context"
	"proto/scyllapb"
	"time"
)

func UploadMessages(client scyllapb.ScyllaServiceClient, messages []*scyllapb.ChatMessage) (*scyllapb.BaseResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.UploadMessages(ctx, &scyllapb.UploadMessagesRequest{
		Messages: messages,
	})
}
