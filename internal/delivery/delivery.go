package delivery

import "fmt"

// DigestSender defines the interface for sending digests.
type DigestSender interface {
	SendDigest(chatID int64, digest string) error
}

// TelegramDigestSender is the production implementation using Telegram.
type TelegramDigestSender struct {
	// Add fields needed for Telegram client/bot, e.g., botToken string, client *tgbotapi.BotAPI
}

// NewTelegramDigestSender creates a new instance of TelegramDigestSender.
func NewTelegramDigestSender() *TelegramDigestSender {
	// TODO: Initialize Telegram client/bot here
	return &TelegramDigestSender{}
}

// SendDigest implements the DigestSender interface.
func (s *TelegramDigestSender) SendDigest(chatID int64, digest string) error {
	// TODO: Implement sending digest via Telegram
	fmt.Printf("TelegramDigestSender: SendDigest called for chat %d (not implemented)\n", chatID)
	return nil
}
