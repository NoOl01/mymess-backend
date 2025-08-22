package db

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	Id        gocql.UUID `db:"message_id"`
	ChatId    gocql.UUID `db:"chat_id"`
	SenderId  int64      `db:"sender_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
	Edited    bool       `db:"edited"`
	Deleted   bool       `db:"deleted"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type PrivateChat struct {
	Id            gocql.UUID  `db:"id"`
	CreatedAt     time.Time   `db:"created_at"`
	LastMessageId *gocql.UUID `db:"last_message_id"`
}

type GroupChat struct {
	Id            gocql.UUID  `db:"id"`
	Title         string      `db:"title"`
	Avatar        string      `db:"avatar"`
	Banner        string      `db:"banner"`
	CreatedAt     time.Time   `db:"created_at"`
	LastMessageId *gocql.UUID `db:"last_message_id"`
}

type ChatParticipant struct {
	ChatId  gocql.UUID `db:"chat_id"`
	UserId  int64      `db:"user_id"`
	Unread  int        `db:"unread"`
	AddedAt *time.Time `db:"added_at"`
}
