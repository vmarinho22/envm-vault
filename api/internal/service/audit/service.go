package audit

import (
	"context"
	"log/slog"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
)

// Logger is the audit API used by middleware and domain services.
type Logger interface {
	AttachOrigin(ctx context.Context, origin RequestOrigin) context.Context
	WithChannel(ctx context.Context, channel string) context.Context
	Record(ctx context.Context, event Event)
}

// Service persists audit events and attaches request origin to context.
type Service struct {
	logs  *repository.AuditLogRepository
	users *repository.UserRepository
}

func New(logs *repository.AuditLogRepository, users *repository.UserRepository) *Service {
	return &Service{logs: logs, users: users}
}

func (s *Service) AttachOrigin(ctx context.Context, origin RequestOrigin) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return withOrigin(ctx, origin)
}

func (s *Service) WithChannel(ctx context.Context, channel string) context.Context {
	origin := OriginFromContext(ctx)
	origin.AuthChannel = channel
	return withOrigin(ctx, origin)
}

// Record writes an audit row and mirrors to slog. Failures are logged and do not fail the request.
func (s *Service) Record(ctx context.Context, event Event) {
	if event.Action == "" {
		return
	}

	origin := OriginFromContext(ctx)
	now := time.Now().UTC()
	row := &model.AuditLog{
		OccurredAt:   now,
		Action:       event.Action,
		ActorUserID:  event.ActorUserID,
		ActorEmail:   event.ActorEmail,
		ActorRole:    event.ActorRole,
		TargetUserID: event.TargetUserID,
		AuthChannel:  origin.AuthChannel,
		IP:           origin.IP,
		UserAgent:    origin.UserAgent,
		Origin:       origin.Origin,
		CreatedAt:    now,
	}

	slog.Info("audit",
		"action", event.Action,
		"auth_channel", origin.AuthChannel,
		"actor_email", event.ActorEmail,
		"actor_role", event.ActorRole,
		"ip", origin.IP,
		"origin", origin.Origin,
	)

	if s == nil || s.logs == nil {
		return
	}
	if err := s.logs.Create(ctx, row); err != nil {
		slog.Error("audit log persist failed", "action", event.Action, "error", err)
	}
}

// Nop is a no-op Logger for tests.
type Nop struct{}

// OrNop returns Nop when l is nil so callers can store a non-nil Logger.
func OrNop(l Logger) Logger {
	if l == nil {
		return Nop{}
	}
	return l
}

func (Nop) AttachOrigin(ctx context.Context, origin RequestOrigin) context.Context {
	return withOrigin(ctx, origin)
}

func (Nop) WithChannel(ctx context.Context, channel string) context.Context {
	origin := OriginFromContext(ctx)
	origin.AuthChannel = channel
	return withOrigin(ctx, origin)
}

func (Nop) Record(context.Context, Event) {}

// Spy is a Logger that records Events for tests.
type Spy struct {
	Events  []Event
	Origins []RequestOrigin
}

func (s *Spy) AttachOrigin(ctx context.Context, origin RequestOrigin) context.Context {
	s.Origins = append(s.Origins, origin)
	return Nop{}.AttachOrigin(ctx, origin)
}

func (s *Spy) WithChannel(ctx context.Context, channel string) context.Context {
	return Nop{}.WithChannel(ctx, channel)
}

func (s *Spy) Record(_ context.Context, event Event) {
	s.Events = append(s.Events, event)
}
