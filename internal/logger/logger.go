package logger

import (
	"log/slog"
	"os"
	"todoApi/internal/config"
)

const (
	colorReset = "\033[0m"
	colorDebug = "\033[36m"
	colorInfo  = "\033[32m"
	colorWarn  = "\033[33m"
	colorError = "\033[31m"
)

func New(cfg *config.Config) *slog.Logger {
	switch cfg.Env {
	case "local":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "prod":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
