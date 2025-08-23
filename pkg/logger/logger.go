package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type keyType int

const loggerKey = keyType(0)

func InitLogger(cfgLevel string) error {
	var (
		lvl       slog.Level
		addSource bool
	)

	switch strings.ToLower(cfgLevel) {
	case "error":
		lvl = slog.LevelError
	case "warn":
		lvl = slog.LevelWarn
	case "info":
		lvl = slog.LevelInfo
	case "debug":
		lvl = slog.LevelDebug
		addSource = true
	default:
		return fmt.Errorf("invalid logger level: %s", cfgLevel)
	}

	handler := slog.Handler(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: addSource,
			Level:     lvl,
		},
	))

	slog.SetDefault(slog.New(handler))

	return nil
}

func FromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey)
	if v == nil {
		return slog.Default()
	}

	logger := v.(*slog.Logger)
	return logger
}

func ContextWithSlogLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func InfoAddValues(ctx context.Context, message string, args ...any) context.Context {
	l := FromContext(ctx)
	l = l.With(args...)
	l.Info(message)
	return ContextWithSlogLogger(ctx, l)
}

func Info(ctx context.Context, message string, args ...any) {
	l := FromContext(ctx)
	l.Info(message, args...)
}

func Error(ctx context.Context, message string, args ...any) {
	l := FromContext(ctx)
	l.Error(message, args...)
}

func Debug(ctx context.Context, message string, args ...any) {
	l := FromContext(ctx)
	l.Debug(message, args...)
}

func AddValuesToContext(ctx context.Context, args ...any) context.Context {
	l := FromContext(ctx)
	l = l.With(args...)
	return ContextWithSlogLogger(ctx, l)
}

func Fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
