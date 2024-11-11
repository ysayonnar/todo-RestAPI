package logger

import "log/slog"

func New() *slog.Logger {
	return &slog.Logger{}
}
