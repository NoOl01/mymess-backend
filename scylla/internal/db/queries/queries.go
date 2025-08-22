package queries

// Update
var (
	UpdateUnreadQuery                 = "UPDATE chat_participant SET unread = ? WHERE chat_id = ? AND user_id = ?"
	UpdateGroupChatSettingsQuery      = "UPDATE group_chat SET title = ?, avatar = ?, banner = ? WHERE id = ?"
	UpdateGroupChatLastMessageQuery   = "UPDATE group_chat SET last_message_id = ? WHERE id = ?"
	UpdatePrivateChatLastMessageQuery = "UPDATE private_chat SET last_message_id = ? WHERE id = ?"
	UpdateMessageQuery                = "UPDATE message SET content = ?, edited = ? WHERE id = ?"
	UpdateDeleteMessageQuery          = "UPDATE message SET deleted = ?, deleted_at = ? WHERE id = ?"
)

// Get
var (
	GetChatsByUserIdQuery   = "SELECT chat_id FROM chat_participant WHERE user_id = ?"
	GetPrivateChatByIdQuery = "SELECT * FROM private_chat WHERE chat_id IN = ?"
	GetGroupChatByIdQuery   = "SELECT * FROM group_chat where chat_id IN = ?"
)
