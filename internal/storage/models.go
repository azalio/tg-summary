package storage

import (
)

// Chat — информация о чате/группе/супергруппе
type Chat struct {
	ID    int64  `gorm:"primaryKey;column:id"`
	Title string `gorm:"not null"`
	Type  string `gorm:"not null"` // group, supergroup, private, etc.
}

// User — информация об авторе сообщения
type User struct {
	ID          int64  `gorm:"primaryKey;column:id"`
	Username    string
	DisplayName string
}

// Message — сообщение, ссылающееся на чат и пользователя
type Message struct {
	ID                int64      `gorm:"primaryKey;autoIncrement"`
	ChatID            int64      `gorm:"not null;index:idx_messages_chat_time,priority:1"`
	MessageID         int64      `gorm:"not null;uniqueIndex:idx_chat_message"`
	AuthorID          int64      `gorm:"index"`
	Text              string
	Timestamp         int64      `gorm:"not null;index:idx_messages_chat_time,priority:2"`
	ReplyToMessageID  *int64     `gorm:"index"`
	Chat              Chat       `gorm:"foreignKey:ChatID;references:ID"`
	Author            User       `gorm:"foreignKey:AuthorID;references:ID"`
}

// TableName overrides for GORM pluralization
func (Chat) TableName() string    { return "chats" }
func (User) TableName() string    { return "users" }
func (Message) TableName() string { return "messages" }