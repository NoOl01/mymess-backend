package queries

import (
	"fmt"
	"github.com/gocql/gocql"
	"results/errs"
	"scylla_db/internal/db"
	"strings"
	"time"
)

type MessageStruct struct {
	ChatId  string
	UserId  int64
	Message string
}

func InsertPrivateChat(userIds []int64) error {
	if len(userIds) == 0 {
		return fmt.Errorf("%v\n", errs.InvalidRequestBody)
	}

	chatStmt, err := db.PrivateChatTable.Insert()
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}

	participantStmt, err := db.ChatParticipantTable.Insert()
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}

	batch := db.Session.NewBatch(gocql.UnloggedBatch)
	chatId := gocql.TimeUUID()

	batch.Query(
		chatStmt,
		chatId,
		time.Now(),
		nil,
	)

	for _, userId := range userIds {
		batch.Query(
			participantStmt,
			chatId,
			userId,
			0,
			nil,
		)
	}

	if err := db.Session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("execute batch: %v", err)
	}

	return nil
}

func InsertGroupChat(title, avatar, banner string) error {
	chat := db.GroupChat{
		Id:            gocql.TimeUUID(),
		Title:         title,
		Avatar:        avatar,
		Banner:        banner,
		CreatedAt:     time.Now(),
		LastMessageId: nil,
	}

	stmt, _ := db.GroupChatTable.Insert()
	query := db.Session.Query(stmt)

	values := []interface{}{
		chat.Id,
		chat.Title,
		chat.Avatar,
		chat.Banner,
		chat.CreatedAt,
		chat.LastMessageId,
	}

	if err := query.Bind(values...).Exec(); err != nil {
		return err
	}

	return nil
}

func InsertChatParticipant(chatId gocql.UUID, userId int64) error {
	now := time.Now()
	chatParticipant := db.ChatParticipant{
		ChatId:  chatId,
		UserId:  userId,
		Unread:  0,
		AddedAt: &now,
	}

	stmt, _ := db.ChatParticipantTable.Insert()
	query := db.Session.Query(stmt)

	values := []interface{}{
		chatId,
		userId,
		chatParticipant.AddedAt,
		chatParticipant.Unread,
	}

	if err := query.Bind(values...).Exec(); err != nil {
		return err
	}

	return nil
}

func InsertPrivateChatBatch(userIds, interlocutorIds []int64) (map[int64]string, error) {
	chatStmt := "INSERT INTO private_chats (id, created_at, last_message_id) VALUES (?, ?, ?)"
	participantStmt := "INSERT INTO chat_participants (chat_id, user_id, unread, added_at) VALUES (?, ?, ?, ?)"

	chats := make(map[int64]string)
	batch := db.Session.NewBatch(gocql.UnloggedBatch)

	for i := 0; i < len(userIds); i++ {
		chatId := gocql.TimeUUID()
		chatIdStr := chatId.String()

		chats[userIds[i]] = chatIdStr
		chats[interlocutorIds[i]] = chatIdStr

		batch.Query(chatStmt, chatId, time.Now(), nil)

		for _, userId := range []int64{userIds[i], interlocutorIds[i]} {
			batch.Query(participantStmt, chatId, userId, 0, time.Now())
		}
	}

	if err := db.Session.ExecuteBatch(batch); err != nil {
		return nil, fmt.Errorf("failed to insert private chat batch: %v", err)
	}

	return chats, nil
}

func InsertMessageBatch(messages []MessageStruct) error {
	if len(messages) == 0 {
		return errs.MessagesListEmpty
	}

	stmt := "INSERT INTO messages (message_id, chat_id, sender_id, content, created_at, edited, deleted, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	batch := db.Session.NewBatch(gocql.UnloggedBatch)

	for _, message := range messages {
		cleanChatId := strings.TrimSpace(strings.Trim(message.ChatId, `"`))
		chatUUID, err := gocql.ParseUUID(cleanChatId)
		if err != nil {
			return fmt.Errorf("invalid chat UUID '%s': %v", message.ChatId, err)
		}

		batch.Query(
			stmt,
			gocql.TimeUUID(),
			chatUUID,
			message.UserId,
			message.Message,
			time.Now(),
			false,
			false,
			nil,
		)
	}

	if err := db.Session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("batch execution failed: %v", err)
	}

	return nil
}
