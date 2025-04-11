package main

import (
	"context"
	"log" // Standard logger only for initial fatal error during logger setup

	"go.uber.org/zap"

	"github.com/azalio/tg-summary/internal/config"
	"github.com/azalio/tg-summary/internal/delivery"
	applog "github.com/azalio/tg-summary/internal/log"
	"github.com/azalio/tg-summary/internal/scheduler"
	"github.com/azalio/tg-summary/internal/storage"
	"github.com/azalio/tg-summary/internal/summarizer"
	"github.com/azalio/tg-summary/internal/telegram"
	telegramtd "github.com/gotd/td/telegram"
)

func main() {
	// --- Initialize Logger ---
	logger, cleanup := initLogger() // Initialize our logger interface
	defer cleanup()                 // Ensure logs are flushed on exit
	logger.Info("tg-summary service starting...")

	// --- Initialize components (using production stubs) ---
	logger.Info("Initializing components...")

	// --- Load config ---
	cfg, err := config.Load(logger.Named("config"))
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// --- Pass config to components ---
	tgClient, err := telegram.NewRealTelegramClient(logger.Named("telegram"), cfg)
	if err != nil {
		logger.Fatal("Failed to initialize Telegram client", zap.Error(err))
	}
	msgStorage := storage.NewSqliteStorage()           // TODO: Pass logger/config if needed
	llmSummarizer := summarizer.NewOpenAISummarizer()  // TODO: Pass logger/config if needed
	digestSender := delivery.NewTelegramDigestSender() // TODO: Pass logger/config if needed
	taskScheduler := scheduler.NewCronScheduler()      // TODO: Pass logger/config if needed

	// --- Basic check of component linkage ---
	ctx := context.Background()

	// Telegram Client Check + бизнес-логика внутри client.Run
	err = tgClient.Run(ctx, func(ctx context.Context, client *telegramtd.Client) error {
		logger.Info("Telegram client authorized (session is alive)")

		// Получить список групп
		groups, err := tgClient.ListGroups(ctx, client)
		if err != nil {
			logger.Error("Failed to list groups after authorization", zap.Error(err))
		} else {
			for _, g := range groups {
				logger.Info("Group found",
					zap.Int64("chat_id", g.ChatID),
					zap.String("title", g.Title),
					zap.String("type", string(g.Type)),
				)
			}
			logger.Info("ListGroups succeeded", zap.Int("group_count", len(groups)))
		}

		// Здесь можно вызывать другие бизнес-методы (FetchMessages, Summarize, Delivery и т.д.)
		return nil
	})
	if err != nil {
		logger.Fatal("Telegram client run failed", zap.Error(err))
	}

	// Storage Check (Example: Save a dummy message)
	dummyMsg := telegram.Message{ID: 1, ChatID: 123, Text: "test"}
	err = msgStorage.SaveMessage(dummyMsg) // Assign to existing err
	if err != nil {
		logger.Fatal("Failed to save message (stub)", zap.Error(err))
	}
	logger.Info("Message saved to storage (stub).")

	// Summarizer Check (Example: Summarize dummy message)
	summary, err := llmSummarizer.Summarize([]telegram.Message{dummyMsg}) // Re-declare summary, assign to existing err
	if err != nil {
		logger.Fatal("Failed to summarize (stub)", zap.Error(err))
	}
	logger.Info("Summarizer generated summary (stub)", zap.String("summary", summary))

	// Delivery Check (Example: Send dummy summary)
	err = digestSender.SendDigest(123, summary) // Assign to existing err
	if err != nil {
		logger.Fatal("Failed to send digest (stub)", zap.Error(err))
	}
	logger.Info("Digest sent (stub).")

	// Scheduler Check
	err = taskScheduler.Start() // Assign to existing err
	if err != nil {
		logger.Fatal("Failed to start scheduler (stub)", zap.Error(err))
	}
	logger.Info("Scheduler started (stub).")

	// --- Keep service running (example) ---
	// In a real app, you'd likely wait for a signal or the scheduler to finish.
	// For now, just print completion.
	logger.Info("tg-summary service initialization check complete.")

	// Example: Stop scheduler after check
	err = taskScheduler.Stop() // Assign to existing err
	if err != nil {
		logger.Warn("Failed to stop scheduler (stub)", zap.Error(err))
	}
	logger.Info("Scheduler stopped (stub).")
}

// initLogger initializes the application logger and returns it along with a cleanup function.
func initLogger() (applog.Logger, func()) {
	logger, cleanup, err := applog.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logger, cleanup
}
