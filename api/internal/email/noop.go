package email

import (
	"context"
	"fmt"
)

type NoopSender struct{}

func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

func (s *NoopSender) SendLoginCode(_ context.Context, _ string, _ string) error {
	return fmt.Errorf("email sender not configured for production; set APP_ENV=development or configure SMTP")
}
