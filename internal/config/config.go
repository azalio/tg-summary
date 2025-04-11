package config

import (
	"os"
	"strconv"

	applog "github.com/azalio/tg-summary/internal/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config holds all application configuration.
type Config struct {
	TelegramAppID      int
	TelegramAppHash    string
	TelegramPhone      string
	TelegramSessionDir string
	SqlitePath         string // путь до файла SQLite
	// Add other config fields as needed
}

// Load loads configuration from .env (if present) and environment variables.
// All errors and warnings are logged using the provided logger.
func Load(logger applog.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Debug(".env file not found or failed to load", zap.Error(err))
	}

	appIDStr := os.Getenv("TELEGRAM_APP_ID")
	appHash := os.Getenv("TELEGRAM_APP_HASH")
	phone := os.Getenv("TELEGRAM_PHONE")
	sessionDir := os.Getenv("TELEGRAM_SESSION_DIR")
	sqlitePath := os.Getenv("SQLITE_PATH")

	missing := false
	if appIDStr == "" {
		logger.Error("Missing TELEGRAM_APP_ID in environment")
		missing = true
	}
	if appHash == "" {
		logger.Error("Missing TELEGRAM_APP_HASH in environment")
		missing = true
	}
	if phone == "" {
		logger.Error("Missing TELEGRAM_PHONE in environment")
		missing = true
	}
	if sessionDir == "" {
		logger.Error("Missing TELEGRAM_SESSION_DIR in environment")
		missing = true
	}
	if sqlitePath == "" {
		logger.Error("Missing SQLITE_PATH in environment")
		missing = true
	}
	if missing {
		return nil, ErrMissingConfig
	}

	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		logger.Error("Invalid TELEGRAM_APP_ID, must be integer", zap.String("value", appIDStr), zap.Error(err))
		return nil, err
	}

	return &Config{
		TelegramAppID:      appID,
		TelegramAppHash:    appHash,
		TelegramPhone:      phone,
		TelegramSessionDir: sessionDir,
		SqlitePath:         sqlitePath,
	}, nil
}

// ErrMissingConfig is returned when required config is missing.
var ErrMissingConfig = &ConfigError{"missing required Telegram config in environment"}

// ConfigError represents a configuration error.
type ConfigError struct {
	Msg string
}

func (e *ConfigError) Error() string { return e.Msg }