package db_models

import "time"

type User struct {
	Id           int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nickname     string      `gorm:"not null;size:30" json:"nickname"`
	Username     string      `gorm:"size:30" json:"username"`
	Email        string      `gorm:"unique;not null;size:50" json:"email"`
	Password     string      `gorm:"not null;" json:"password"`
	AvatarPath   string      `gorm:"size:255" json:"avatar_path"`
	BannerPath   string      `gorm:"size:255" json:"banner_path"`
	RegisterDate time.Time   `json:"register_date"`
	Deleted      bool        `json:"deleted"`
	DeletedAt    time.Time   `json:"deleted_at"`
	UserThemes   []UserTheme `gorm:"foreignKey:UserId" json:"user_themes"`
}

type Chat struct {
	UserId        int64     `gorm:"not null;index:idx_user_peer" json:"user_id"`
	User          *User     `gorm:"foreignKey:UserId" json:"-"`
	PeerId        int64     `gorm:"not null;index:idx_user_peer" json:"peer_id"`
	Peer          *User     `gorm:"foreignKey:PeerId" json:"-"`
	LastMessageId int64     `json:"last_message_id"`
	LastMessage   *Message  `gorm:"foreignKey:LastMessageId" json:"-"`
	UnreadCount   int       `gorm:"not null;default:0" json:"unread_count"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Message struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"not null;" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserId" json:"-"`
	SenderId  int64     `gorm:"not null;" json:"sender_id"`
	Sender    *User     `gorm:"foreignKey:SenderId" json:"-"`
	Content   string    `gorm:"not null;" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Edited    bool      `json:"edited"`
}

type UserTheme struct {
	Id                      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId                  int64  `gorm:"not null;" json:"user_id"`
	User                    User   `gorm:"foreignKey:UserId" json:"user"`
	ChatBackgroundImagePath string `gorm:"size:255" json:"chat_background_image_path"`
}
