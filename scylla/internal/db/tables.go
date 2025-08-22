package db

import "github.com/scylladb/gocqlx/table"

var MessageTable = table.New(table.Metadata{
	Name: "message",
	Columns: []string{
		"id", "chat_id", "sender_id", "content", "created_at", "edited", "deleted", "deleted_at",
	},
	PartKey: []string{"id"},
	SortKey: []string{"created_at"},
})

var PrivateChatTable = table.New(table.Metadata{
	Name: "private_chat",
	Columns: []string{
		"id", "user_id", "interlocutor_id", "unread", "created_at", "last_message_id",
	},
	PartKey: []string{"id"},
})

var GroupChatTable = table.New(table.Metadata{
	Name: "group_chat",
	Columns: []string{
		"id", "title", "avatar", "banner", "created_at", "last_message_id",
	},
	PartKey: []string{"id"},
	SortKey: []string{"created_at"},
})

var ChatParticipantTable = table.New(table.Metadata{
	Name: "chat_participant",
	Columns: []string{
		"chat_id", "user_id", "unread", "added_at",
	},
	PartKey: []string{"chat_id"},
})
