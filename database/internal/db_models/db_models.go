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

type DeviceToken struct {
	Id     int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId int64  `gorm:"not null;" json:"user_id"`
	User   User   `gorm:"foreignKey:UserId" json:"-"`
	Token  string `gorm:"not null" json:"token"`
}

type UserTheme struct {
	Id                              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId                          int64  `gorm:"not null;" json:"user_id"`
	User                            User   `gorm:"foreignKey:UserId" json:"user"`
	ThemeName                       string `gorm:"size:60" json:"theme_name"`
	ChatBackgroundImagePath         string `gorm:"size:255" json:"chat_background_image_path"`
	AccentColor                     string `gorm:"size:50" json:"accent_color"`
	BackgroundColor                 string `gorm:"size:50" json:"background_color"`
	ChatUserBackgroundColor         string `gorm:"size:50" json:"chat_user_background_color"`
	ChatInterlocutorBackgroundColor string `gorm:"size:50" json:"chat_interlocutor_background_color"`
	TextColor                       string `gorm:"size:50" json:"text_color"`
}
