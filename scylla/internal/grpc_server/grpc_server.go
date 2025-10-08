package grpc_server

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "proto/scyllapb"
	"results/succ"
	"scylla_db/internal/db/queries"
	"strconv"
	"strings"
	"time"
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

func (s *Server) GetChats(_ context.Context, req *pb.UserId) (*pb.ChatsResponse, error) {
	userId := req.UserId

	chats, err := queries.GetUserChats(userId)
	if err != nil {
		return &pb.ChatsResponse{
			Error: err.Error(),
		}, nil
	}

	var pbChats []*pb.Chats

	for _, chat := range chats {
		pbChat := &pb.Chats{}

		if chatId, ok := chat["chat_id"].(gocql.UUID); ok {
			pbChat.ChatId = chatId.String()
		}

		if lastMessageId, ok := chat["last_message_id"].(*gocql.UUID); ok && lastMessageId != nil {
			message, err := queries.GetMessageById(*lastMessageId)
			if err == nil {
				pbChat.LastMessage = &pb.Message{
					MessageId: message.Id.String(),
					ChatId:    message.ChatId.String(),
					SenderId:  message.SenderId,
					Content:   message.Content,
					CreatedAt: message.CreatedAt.Format(time.RFC3339),
					Edited:    message.Edited,
					Deleted:   message.Deleted,
				}
			}
		}

		if participants, ok := chat["participants"].([]int64); ok {
			pbChat.Participants = participants
		}

		pbChats = append(pbChats, pbChat)
	}

	return &pb.ChatsResponse{
		Chats: pbChats,
	}, nil
}

func (s *Server) GetChatHistory(_ context.Context, req *pb.ChatId) (*pb.ChatHistoryResponse, error) {
	chatIdStr := req.ChatId
	chatId, err := gocql.ParseUUID(chatIdStr)
	if err != nil {
		return &pb.ChatHistoryResponse{
			Error: err.Error(),
		}, err
	}

	messages, err := queries.GetChatHistory(chatId)
	if err != nil {
		return &pb.ChatHistoryResponse{
			Error: err.Error(),
		}, err
	}

	var pbMessages []*pb.Message
	for _, m := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			MessageId: m.Id.String(),
			ChatId:    m.ChatId.String(),
			SenderId:  m.SenderId,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
			Edited:    m.Edited,
			Deleted:   m.Deleted,
		})
	}

	return &pb.ChatHistoryResponse{
		Messages: pbMessages,
	}, nil
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
