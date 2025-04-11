package log

import "go.uber.org/zap"

// Logger defines a standard interface for logging.
// This allows swapping the underlying logging implementation.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Named(name string) Logger // Method to create a named sub-logger
}

// zapLogger is an adapter for zap.Logger to satisfy the Logger interface.
// zapLogger is an adapter for zap.Logger to satisfy the Logger interface.
type ZapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new Logger implementation using zap.
// In the future, this can be extended to support other backends.
func NewLogger() (Logger, func(), error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, func(){}, err
	}
	cleanup := func() { _ = zapLogger.Sync() }
	return &ZapLogger{logger: zapLogger}, cleanup, nil
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Named creates a new named logger derived from the current one.
func (l *ZapLogger) Named(name string) Logger {
	return &ZapLogger{logger: l.logger.Named(name)}
}