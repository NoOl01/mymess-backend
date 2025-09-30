package queries

import (
	"github.com/gocql/gocql"
	"scylla_db/internal/db"
)

func GetUserChats(userId int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	query := "SELECT chat_id FROM chat_participants WHERE user_id = ?"
	iter := db.Session.Query(query, userId).Iter()

	var chatId gocql.UUID
	var chatIds []gocql.UUID

	for iter.Scan(&chatId) {
		chatIds = append(chatIds, chatId)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	for _, chatId := range chatIds {
		chatInfo := make(map[string]interface{})
		chatInfo["chat_id"] = chatId

		var lastMessageId *gocql.UUID
		err := db.Session.Query("SELECT last_message_id FROM private_chats WHERE id = ?", chatId).Scan(&lastMessageId)

		if err == nil {
			chatInfo["last_message_id"] = lastMessageId

			var participants []int64
			participantsQuery := "SELECT user_id FROM chat_participants WHERE chat_id = ?"
			participantsIter := db.Session.Query(participantsQuery, chatId).Iter()

			var participantId int64
			for participantsIter.Scan(&participantId) {
				if participantId != userId {
					participants = append(participants, participantId)
				}
			}
			participantsIter.Close()

			chatInfo["participants"] = participants
		} else {
			err = db.Session.Query("SELECT last_message_id FROM group_chats WHERE id = ?", chatId).Scan(&lastMessageId)

			if err == nil {
				chatInfo["last_message_id"] = lastMessageId

				var participants []int64
				participantsQuery := "SELECT user_id FROM chat_participants WHERE chat_id = ?"
				participantsIter := db.Session.Query(participantsQuery, chatId).Iter()

				var participantId int64
				for participantsIter.Scan(&participantId) {
					participants = append(participants, participantId)
				}
				participantsIter.Close()

				chatInfo["participants"] = participants
			} else {
				continue
			}
		}

		result = append(result, chatInfo)
	}

	return result, nil
}

func GetMessageById(messageId gocql.UUID) (*db.Message, error) {
	query := "SELECT * FROM messages WHERE message_id = ?"

	var message db.Message
	err := db.Session.Query(query, messageId).Scan(
		&message.Id,
		&message.ChatId,
		&message.SenderId,
		&message.Content,
		&message.CreatedAt,
		&message.Edited,
		&message.Deleted,
		&message.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &message, nil
}
