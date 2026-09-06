package audit_test

import (
	"context"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
)

func TestAttachOrigin_RoundTrip(t *testing.T) {
	svc := audit.New(nil)
	ctx := svc.AttachOrigin(context.Background(), audit.RequestOrigin{
		IP:          "203.0.113.10",
		UserAgent:   "vitest",
		Origin:      "http://localhost:3100",
		AuthChannel: audit.ChannelLoginAttempt,
	})

	got := audit.OriginFromContext(ctx)
	if got.IP != "203.0.113.10" || got.UserAgent != "vitest" || got.Origin != "http://localhost:3100" {
		t.Fatalf("origin = %#v", got)
	}
	if got.AuthChannel != audit.ChannelLoginAttempt {
		t.Fatalf("channel = %q", got.AuthChannel)
	}

	ctx = svc.WithChannel(ctx, audit.ChannelLoggedInUser)
	if audit.OriginFromContext(ctx).AuthChannel != audit.ChannelLoggedInUser {
		t.Fatal("expected logged_in_user channel after WithChannel")
	}
}

func TestRecord_DoesNotPanicWithoutRepository(t *testing.T) {
	svc := audit.New(nil)
	actor := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	svc.Record(context.Background(), audit.Event{
		Action:      audit.ActionPersonalTokenCreated,
		Summary:     audit.SummaryPersonalTokenCreated,
		ActorUserID: &actor,
		ActorEmail:  "user@example.com",
	})
}

func TestNopAndSpy_ImplementLogger(t *testing.T) {
	var _ audit.Logger = audit.Nop{}
	var _ audit.Logger = &audit.Spy{}
}

func TestOrNop(t *testing.T) {
	if _, ok := audit.OrNop(nil).(audit.Nop); !ok {
		t.Fatal("OrNop(nil) should return Nop")
	}
	svc := audit.New(nil)
	if audit.OrNop(svc) != svc {
		t.Fatal("OrNop should return the same non-nil logger")
	}
}

func TestSpy_Record(t *testing.T) {
	spy := &audit.Spy{}
	spy.Record(context.Background(), audit.Event{Action: audit.ActionLoginFailed})
	if len(spy.Events) != 1 || spy.Events[0].Action != audit.ActionLoginFailed {
		t.Fatalf("events = %#v", spy.Events)
	}
}
