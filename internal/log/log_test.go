package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerInterface(t *testing.T) {
	logger, cleanup, err := NewLogger()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer cleanup()

	// Test all interface methods (should not panic)
	logger.Debug("debug message", zap.String("key", "value"))
	logger.Info("info message", zap.Int("int", 42))
	logger.Warn("warn message", zap.Float64("float", 3.14))
	logger.Error("error message", zap.Bool("bool", true))
	// Fatal is not called in tests as it would exit the process

	// Test Named logger
	sub := logger.Named("subsystem")
	sub.Info("named logger works", zapcore.Field{Key: "test", Type: zapcore.StringType, String: "ok"})
}