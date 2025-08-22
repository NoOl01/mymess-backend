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

type newChats struct {
	UserId         int64
	InterlocutorId int64
	Content        string
}

func (s *Server) UploadMessages(_ context.Context, req *pb.UploadMessagesRequest) (*pb.BaseResultResponse, error) {
	messages := req.GetMessages()
	pendingMessages := make(map[string][]newChats)

	for _, message := range messages {
		if strings.HasPrefix(message.ChatId, "-1") && message.ChatId != "" {
			interlocutorIdStr := message.ChatId[2:]
			interlocutorId, err := strconv.ParseInt(interlocutorIdStr, 10, 64)
			if err != nil {
				return &pb.BaseResultResponse{
					Result: fmt.Sprintf("Invalid interlocutor ID: %s", err.Error()),
				}, err
			}
			pendingMessages[message.ChatId] = append(pendingMessages[message.ChatId], newChats{
				UserId:         message.UserId,
				InterlocutorId: interlocutorId,
				Content:        message.Message,
			})
		} else {
			continue
		}
	}

	var chats map[int64]string
	if len(pendingMessages) > 0 {
		var userIds []int64
		var interlocutorIds []int64
		uniquePairs := make(map[string]bool)

		for _, chatList := range pendingMessages {
			for _, chat := range chatList {
				key := generateChatKey(chat.UserId, chat.InterlocutorId)
				if uniquePairs[key] {
					continue
				}
				uniquePairs[key] = true

				userIds = append(userIds, chat.UserId)
				interlocutorIds = append(interlocutorIds, chat.InterlocutorId)
			}
		}

		var err error
		chats, err = queries.InsertPrivateChatBatch(userIds, interlocutorIds)
		if err != nil {
			return &pb.BaseResultResponse{
				Result: err.Error(),
			}, err
		}
	}

	var newMessages []queries.MessageStruct

	for _, message := range messages {
		if strings.HasPrefix(message.ChatId, "-1") {
			chatUUID := chats[message.UserId]
			if chatUUID == "" {
				return &pb.BaseResultResponse{
					Result: "Chat UUID not found for new chat",
				}, fmt.Errorf("chat UUID not found for user %d", message.UserId)
			}

			newMessages = append(newMessages, queries.MessageStruct{
				ChatId:  chatUUID,
				UserId:  message.UserId,
				Message: message.Message,
			})
		} else {
			newMessages = append(newMessages, queries.MessageStruct{
				ChatId:  message.ChatId,
				UserId:  message.UserId,
				Message: message.Message,
			})
		}
	}

	if err := queries.InsertMessageBatch(newMessages); err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}

func generateChatKey(id1, id2 int64) string {
	if id1 < id2 {
		return fmt.Sprintf("%d:%d", id1, id2)
	}
	return fmt.Sprintf("%d:%d", id2, id1)
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
