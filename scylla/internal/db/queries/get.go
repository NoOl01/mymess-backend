package queries

import (
	"fmt"
	"github.com/gocql/gocql"
	"scylla_db/internal/db"
)

func GetChatsByUserId(userId int64) ([]gocql.UUID, error) {
	iter := db.Session.Query(GetChatsByUserIdQuery, userId).Iter()
	defer func(iter *gocql.Iter) {
		if err := iter.Close(); err != nil {
			fmt.Println(err)
		}
	}(iter)

	var chatId gocql.UUID
	chatIds := make([]gocql.UUID, 0, 1000)

	for iter.Scan(&chatId) {
		chatIds = append(chatIds, chatId)
	}

	return chatIds, nil
}

func GetPrivatesChatById(chatId []gocql.UUID) ([]db.PrivateChat, error) {
	iter := db.Session.Query(GetPrivateChatByIdQuery, chatId).Iter()
	defer func(iter *gocql.Iter) {
		if err := iter.Close(); err != nil {
			fmt.Println(err)
		}
	}(iter)

	chats := make([]db.PrivateChat, 0, 1000)

	for {
		var chat db.PrivateChat
		if !iter.Scan(&chat) {
			break
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func GetGroupChatsById(chatId []gocql.UUID) ([]db.GroupChat, error) {
	iter := db.Session.Query(GetGroupChatByIdQuery, chatId).Iter()
	defer func(iter *gocql.Iter) {
		if err := iter.Close(); err != nil {
			fmt.Println(err)
		}
	}(iter)

	chats := make([]db.GroupChat, 0, 1000)

	for {
		var chat db.GroupChat
		if !iter.Scan(&chat) {
			break
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
