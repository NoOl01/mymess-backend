package client

import (
	"context"
	"proto/scyllapb"
	"time"
)

func GetChats(client scyllapb.ScyllaServiceClient, userId int64) (*scyllapb.ChatsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.GetChats(ctx, &scyllapb.UserId{
		UserId: userId,
	})
}

func GetChatHistory(client scyllapb.ScyllaServiceClient, chatId string) (*scyllapb.ChatHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.GetChatHistory(ctx, &scyllapb.ChatId{
		ChatId: chatId,
	})
}
