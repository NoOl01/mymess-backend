package queries

import (
	"github.com/gocql/gocql"
	"scylla_db/internal/db"
	"time"
)

func UpdateUnreadCount(chatId gocql.UUID, userId int64, unread int) error {
	if err := db.Session.Query(UpdateUnreadQuery, unread, chatId, userId).Exec(); err != nil {
		return err
	}

	return nil
}

func UpdateGroupChatSettings(id gocql.UUID, title, avatar, banner string) error {
	if err := db.Session.Query(UpdateGroupChatSettingsQuery, title, avatar, banner, id).Exec(); err != nil {
		return err
	}

	return nil
}

func UpdateLastMessageId(id, messageId gocql.UUID, isGroupChat bool) error {
	switch isGroupChat {
	case true:
		if err := db.Session.Query(UpdateGroupChatLastMessageQuery, messageId, id).Exec(); err != nil {
			return err
		}
	case false:
		if err := db.Session.Query(UpdatePrivateChatLastMessageQuery, messageId, id).Exec(); err != nil {
			return err
		}
	}

	return nil
}

func UpdateMessage(id gocql.UUID, content string) error {
	if err := db.Session.Query(UpdateMessageQuery, content, true, id).Exec(); err != nil {
		return err
	}

	return nil
}

func UpdateDeleteMessage(id gocql.UUID) error {
	if err := db.Session.Query(UpdateDeleteMessageQuery, time.Now(), id).Exec(); err != nil {
		return err
	}

	return nil
}
