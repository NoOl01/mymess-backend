package grpc_server

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "proto/scyllapb"
	"results/succ"
	"scylla_db/internal/db/queries"
	"strconv"
	"strings"
)

type Server struct {
	pb.UnimplementedScyllaServiceServer
}

func (s *Server) UploadMessages(_ context.Context, req *pb.UploadMessagesRequest) (*pb.BaseResultResponse, error) {
	messages := req.GetMessages()
	pendingMessages := make(map[string][]queries.MessageData)

	fmt.Printf("Processing %d messages\n", len(messages))

	for _, message := range messages {
		fmt.Printf("Message: ChatId='%s', UserId=%d, Message='%s'\n", message.ChatId, message.UserId, message.Message)

		if strings.HasPrefix(message.ChatId, "-1") && message.ChatId != "" {
			fmt.Printf("Creating new chat for ChatId: %s\n", message.ChatId)
			interlocutorIdStr := message.ChatId[2:]
			interlocutorId, err := strconv.ParseInt(interlocutorIdStr, 10, 64)
			if err != nil {
				return &pb.BaseResultResponse{
					Result: fmt.Sprintf("Invalid interlocutor ID: %s", err.Error()),
				}, err
			}
			pendingMessages[message.ChatId] = append(pendingMessages[message.ChatId], queries.MessageData{
				UserId:         message.UserId,
				InterlocutorId: interlocutorId,
				Content:        message.Message,
			})
		} else {
			fmt.Printf("Using existing chat: %s\n", message.ChatId)
		}
	}

	var chatIdToUUID map[string]string

	if len(pendingMessages) > 0 {
		fmt.Printf("Creating %d new chats\n", len(pendingMessages))
		var err error
		chatIdToUUID, err = queries.InsertPrivateChatBatch(pendingMessages)
		if err != nil {
			return &pb.BaseResultResponse{
				Result: err.Error(),
			}, err
		}
		fmt.Printf("Created chats: %+v\n", chatIdToUUID)
	}

	var newMessages []queries.MessageStruct

	for _, message := range messages {
		if strings.HasPrefix(message.ChatId, "-1") {
			chatUUID := chatIdToUUID[message.ChatId]
			if chatUUID == "" {
				return &pb.BaseResultResponse{
					Result: "Chat UUID not found for new chat",
				}, fmt.Errorf("chat UUID not found for chatId %s", message.ChatId)
			}

			fmt.Printf("Adding message for new chat - Original: %s, UUID: %s\n", message.ChatId, chatUUID)
			newMessages = append(newMessages, queries.MessageStruct{
				ChatId:  chatUUID,
				UserId:  message.UserId,
				Message: message.Message,
				Time:    message.Time.AsTime(),
			})
		} else {
			fmt.Printf("Adding message for existing chat: %s\n", message.ChatId)
			newMessages = append(newMessages, queries.MessageStruct{
				ChatId:  message.ChatId,
				UserId:  message.UserId,
				Message: message.Message,
				Time:    message.Time.AsTime(),
			})
		}
	}

	fmt.Printf("Inserting %d messages\n", len(newMessages))
	if err := queries.InsertMessageBatch(newMessages); err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
