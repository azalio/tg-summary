package telegram

import (
	"context"
	"fmt" // Keep fmt for now, might be used in other methods

	applog "github.com/azalio/tg-summary/internal/log"
	"github.com/azalio/tg-summary/internal/config"
	"bufio"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type Message struct {
	ID        int64
	ChatID    int64
	Sender    string
	Text      string
	Timestamp int64
}

type GroupType string

const (
	GroupTypeGroup      GroupType = "group"
	GroupTypeSupergroup GroupType = "supergroup"
)

type GroupInfo struct {
	ChatID int64
	Title  string
	Type   GroupType
}

type TelegramClient interface {
	Run(ctx context.Context, fn func(ctx context.Context, api *tg.Client) error) error
	FetchMessages(ctx context.Context, chatID int64, from, to int64) ([]Message, error)
	ListGroups(ctx context.Context, api *tg.Client) ([]GroupInfo, error)
	SendMessage(ctx context.Context, chatID int64, text string) error
}

// RealTelegramClient is the production implementation of TelegramClient
type RealTelegramClient struct {
	client     *telegram.Client
	sessionDir string
	appID      int
	appHash    string
	phone      string
	log        applog.Logger
}

// NewRealTelegramClient creates a new instance of RealTelegramClient using injected config and logger.
func NewRealTelegramClient(logger applog.Logger, cfg *config.Config) (*RealTelegramClient, error) {
	logger.Info("Initializing RealTelegramClient with injected config")
	return &RealTelegramClient{
		log:        logger,
		appID:      cfg.TelegramAppID,
		appHash:    cfg.TelegramAppHash,
		phone:      cfg.TelegramPhone,
		sessionDir: cfg.TelegramSessionDir,
	}, nil
}

// Authorize implements the TelegramClient

func (c *RealTelegramClient) Run(ctx context.Context, fn func(ctx context.Context, api *telegram.Client) error) error {
	c.log.Info("Starting Telegram client session")

	// Prepare session storage
	sessionFile := filepath.Join(c.sessionDir, "session.json")
	storage := &session.FileStorage{Path: sessionFile}

	client := telegram.NewClient(c.appID, c.appHash, telegram.Options{
		SessionStorage: storage,
		Logger:         zap.NewNop(), // gotd expects zap.Logger, but we use our own for app logs
	})

	// Run client and execute user callback
	err := client.Run(ctx, func(ctx context.Context) error {
		c.log.Info("Running gotd/td client session")
		authFlow := auth.NewFlow(
			auth.Constant(
				c.phone,
				"", // password (not used)
				auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
					c.log.Info("Enter the code sent to your Telegram account")
					fmt.Print("Telegram code: ")
					reader := bufio.NewReader(os.Stdin)
					code, err := reader.ReadString('\n')
					if err != nil {
						c.log.Error("Failed to read code from stdin", zap.Error(err))
						return "", err
					}
					return code[:len(code)-1], nil // remove newline
				}),
			),
			auth.SendCodeOptions{},
		)
		if err := client.Auth().IfNecessary(ctx, authFlow); err != nil {
			c.log.Error("Telegram authorization failed", zap.Error(err))
			return err
		}
		c.log.Info("Telegram authorization successful")
		return fn(ctx, client)
	})
	if err != nil {
		c.log.Error("Telegram client run failed", zap.Error(err))
		return err
	}
	return nil
}

// FetchMessages implements the TelegramClient interface
func (c *RealTelegramClient) FetchMessages(ctx context.Context, chatID int64, from, to int64) ([]Message, error) {
	// TODO: Implement real message fetching logic
	fmt.Printf("RealTelegramClient: FetchMessages called for chat %d (not implemented)\n", chatID)
	return []Message{}, nil
}

// ListGroups implements the TelegramClient interface
func (c *RealTelegramClient) ListGroups(ctx context.Context, api *telegram.Client) ([]GroupInfo, error) {
	var groups []GroupInfo

	// Get dialogs (chats, groups, supergroups, channels)
	dialogs, err := api.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetDate:   0,
		OffsetID:     0,
		OffsetPeer:   &tg.InputPeerEmpty{}, // fix: must not be nil
		Limit:        100, // adjust as needed
		Hash:         0,
	})
	if err != nil {
		c.log.Error("Failed to get dialogs", zap.Error(err))
		return nil, err
	}

	// Correct way to extract chats from MessagesDialogsClass
	chats := []tg.ChatClass{}
	switch d := dialogs.(type) {
	case *tg.MessagesDialogs:
		chats = d.Chats
	case *tg.MessagesDialogsSlice:
		chats = d.Chats
	case *tg.MessagesDialogsNotModified:
		// No chats
	}

	for _, peer := range chats {
		switch chat := peer.(type) {
		case *tg.Chat:
			if chat.Deactivated || chat.MigratedTo != nil {
				continue
			}
			if chat.ID != 0 && chat.Title != "" {
				groups = append(groups, GroupInfo{
					ChatID: int64(chat.ID),
					Title:  chat.Title,
					Type:   GroupTypeGroup,
				})
			}
		case *tg.Channel:
			if chat.Megagroup {
				groups = append(groups, GroupInfo{
					ChatID: chat.ID,
					Title:  chat.Title,
					Type:   GroupTypeSupergroup,
				})
			}
		}
	}

	return groups, nil
}

// SendMessage implements the TelegramClient interface
func (c *RealTelegramClient) SendMessage(ctx context.Context, chatID int64, text string) error {
	// TODO: Implement real message sending logic
	fmt.Printf("RealTelegramClient: SendMessage called for chat %d (not implemented)\n", chatID)
	return nil
}
