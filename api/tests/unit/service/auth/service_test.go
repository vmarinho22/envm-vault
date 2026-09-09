package auth_test

import (
	"testing"
	"time"

	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/vmarinho22/envm-vault/api/internal/service/audit"
	"github.com/vmarinho22/envm-vault/api/internal/service/auth"
)

func TestVerifyLogin_InvalidCodeLength_RecordsLoginFailed(t *testing.T) {
	spy := &audit.Spy{}
	svc := auth.New(nil, nil, nil, nil, nil, "test-secret", time.Hour, spy)

	_, err := svc.VerifyLogin(t.Context(), "user@example.com", "123")
	if err != service.ErrInvalidCode {
		t.Fatalf("VerifyLogin() error = %v, want ErrInvalidCode", err)
	}
	if len(spy.Events) != 1 {
		t.Fatalf("audit events = %d, want 1", len(spy.Events))
	}
	if spy.Events[0].Action != audit.ActionLoginFailed {
		t.Errorf("action = %q, want %q", spy.Events[0].Action, audit.ActionLoginFailed)
	}
	if spy.Events[0].ActorEmail != "user@example.com" {
		t.Errorf("actor email = %q, want user@example.com", spy.Events[0].ActorEmail)
	}
}
