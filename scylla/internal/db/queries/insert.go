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

type MessageData struct {
	UserId         int64
	InterlocutorId int64
	Content        string
}

type ChatPair struct {
	UserId         int64
	InterlocutorId int64
}

func InsertPrivateChatBatch(pendingMessages map[string][]MessageData) (map[string]string, error) {
	chatStmt := "INSERT INTO private_chats (id, created_at, last_message_id) VALUES (?, ?, ?)"
	participantStmt := "INSERT INTO chat_participants (chat_id, user_id, unread, added_at) VALUES (?, ?, ?, ?)"

	chatIdToUUID := make(map[string]string)
	pairToUUID := make(map[string]string)

	batch := db.Session.NewBatch(gocql.UnloggedBatch)

	uniquePairs := make(map[string]ChatPair)
	for originalChatId, chatList := range pendingMessages {
		if len(chatList) > 0 {
			firstChat := chatList[0]
			pairKey := generateChatKey(firstChat.UserId, firstChat.InterlocutorId)

			if _, exists := uniquePairs[pairKey]; !exists {
				uniquePairs[pairKey] = ChatPair{
					UserId:         firstChat.UserId,
					InterlocutorId: firstChat.InterlocutorId,
				}
			}

			if _, exists := pairToUUID[pairKey]; !exists {
				chatId := gocql.TimeUUID()
				chatIdStr := chatId.String()
				pairToUUID[pairKey] = chatIdStr

				batch.Query(chatStmt, chatId, time.Now(), nil)

				for _, userId := range []int64{firstChat.UserId, firstChat.InterlocutorId} {
					batch.Query(participantStmt, chatId, userId, 0, time.Now())
				}
			}

			chatIdToUUID[originalChatId] = pairToUUID[pairKey]
		}
	}

	if err := db.Session.ExecuteBatch(batch); err != nil {
		return nil, fmt.Errorf("failed to insert private chat batch: %v", err)
	}

	return chatIdToUUID, nil
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

func generateChatKey(id1, id2 int64) string {
	if id1 < id2 {
		return fmt.Sprintf("%d:%d", id1, id2)
	}
	return fmt.Sprintf("%d:%d", id2, id1)
}
