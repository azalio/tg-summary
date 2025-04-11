package storage

import (
	"context"
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Storage интерфейс для production-уровня
type Storage interface {
	Init(ctx context.Context) error
	SaveChat(ctx context.Context, chat *Chat) error
	SaveUser(ctx context.Context, user *User) error
	SaveMessage(ctx context.Context, msg *Message) error
	GetLastMessageTimestamp(ctx context.Context, chatID int64) (int64, error)
	GetMessagesAfter(ctx context.Context, chatID int64, afterTimestamp int64) ([]Message, error)
	Close() error
}

// GormStorage — production-реализация на GORM
type GormStorage struct {
	db *gorm.DB
}

func NewGormStorage(dsn string) (*GormStorage, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormStorage{db: db}, nil
}

func (s *GormStorage) Init(ctx context.Context) error {
	return s.db.WithContext(ctx).AutoMigrate(&Chat{}, &User{}, &Message{})
}

func (s *GormStorage) SaveChat(ctx context.Context, chat *Chat) error {
	return s.db.WithContext(ctx).Clauses(
		clause.OnConflict{UpdateAll: true},
	).Create(chat).Error
}

func (s *GormStorage) SaveUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Clauses(
		clause.OnConflict{UpdateAll: true},
	).Create(user).Error
}

func (s *GormStorage) SaveMessage(ctx context.Context, msg *Message) error {
	return s.db.WithContext(ctx).Clauses(
		clause.OnConflict{DoNothing: true},
	).Create(msg).Error
}

func (s *GormStorage) GetLastMessageTimestamp(ctx context.Context, chatID int64) (int64, error) {
	var msg Message
	err := s.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("timestamp DESC").
		First(&msg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return msg.Timestamp, nil
}

func (s *GormStorage) GetMessagesAfter(ctx context.Context, chatID int64, afterTimestamp int64) ([]Message, error) {
	var msgs []Message
	err := s.db.WithContext(ctx).
		Where("chat_id = ? AND timestamp > ?", chatID, afterTimestamp).
		Order("timestamp ASC").
		Find(&msgs).Error
	return msgs, err
}

func (s *GormStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
