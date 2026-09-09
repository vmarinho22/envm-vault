package audit_test

import (
	"context"
	"testing"

	"github.com/vmarinho22/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
)

func TestAttachOrigin_RoundTrip(t *testing.T) {
	svc := audit.New(nil, nil)
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
	svc := audit.New(nil, nil)
	actor := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	svc.Record(context.Background(), audit.Event{
		Action:      audit.ActionPersonalTokenCreated,
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
	svc := audit.New(nil, nil)
	if audit.OrNop(svc) != svc {
		t.Fatal("OrNop should return the same non-nil logger")
	}
}

func TestSpy_Record(t *testing.T) {
	spy := &audit.Spy{}
	projectID := uuid.MustParse("00000000-0000-0000-0000-000000000099")
	spy.Record(context.Background(), audit.Event{
		Action:       audit.ActionProjectCreated,
		ProjectID:    &projectID,
		ProjectTitle: "Payment",
		ProjectSlug:  "payment",
	})
	if len(spy.Events) != 1 {
		t.Fatalf("events = %#v", spy.Events)
	}
	ev := spy.Events[0]
	if ev.Action != audit.ActionProjectCreated {
		t.Fatalf("action = %q", ev.Action)
	}
	if ev.ProjectID == nil || *ev.ProjectID != projectID {
		t.Fatalf("project_id = %#v", ev.ProjectID)
	}
	if ev.ProjectTitle != "Payment" || ev.ProjectSlug != "payment" {
		t.Fatalf("project snapshot = %q / %q", ev.ProjectTitle, ev.ProjectSlug)
	}
}
