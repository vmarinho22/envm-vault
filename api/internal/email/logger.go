package email

import (
	"context"
	"log/slog"
)

type LoggerSender struct{}

func NewLoggerSender() *LoggerSender {
	return &LoggerSender{}
}

func (s *LoggerSender) SendLoginCode(_ context.Context, to, code string) error {
	slog.Info("login code sent",
		"email", to,
		"code", code,
		"expires_in", "10m",
	)
	return nil
}
