package token_test

import (
	"strings"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/token"
)

func TestGeneratePlaintext(t *testing.T) {
	const (
		wantPrefix     = "ev_"
		wantRandomHex  = 64
		wantDisplayLen = 11
	)

	first, err := token.GeneratePlaintext()
	if err != nil {
		t.Fatalf("GeneratePlaintext() error = %v", err)
	}
	second, err := token.GeneratePlaintext()
	if err != nil {
		t.Fatalf("GeneratePlaintext() error = %v", err)
	}

	if !strings.HasPrefix(first, wantPrefix) {
		t.Errorf("plaintext = %q, want prefix %q", first, wantPrefix)
	}
	wantLen := len(wantPrefix) + wantRandomHex
	if len(first) != wantLen {
		t.Errorf("len(plaintext) = %d, want %d", len(first), wantLen)
	}
	if first == second {
		t.Error("expected unique plaintext tokens")
	}

	prefix := token.DisplayPrefix(first)
	if prefix != first[:wantDisplayLen] {
		t.Errorf("prefix = %q, want %q", prefix, first[:wantDisplayLen])
	}

	hash := token.HashToken(first)
	if len(hash) != 64 {
		t.Errorf("hash length = %d, want 64", len(hash))
	}
	if hash == token.HashToken(second) {
		t.Error("expected distinct hashes for distinct tokens")
	}
}

func TestCanRevokeTarget(t *testing.T) {
	tests := []struct {
		actor  string
		target string
		want   bool
	}{
		{service.RoleSuperAdmin, service.RoleSuperAdmin, true},
		{service.RoleSuperAdmin, service.RoleAdmin, true},
		{service.RoleSuperAdmin, service.RoleUser, true},
		{service.RoleAdmin, service.RoleUser, true},
		{service.RoleAdmin, service.RoleAdmin, true},
		{service.RoleAdmin, service.RoleSuperAdmin, false},
		{service.RoleUser, service.RoleUser, false},
	}

	for _, tt := range tests {
		got := token.CanRevokeTarget(tt.actor, tt.target)
		if got != tt.want {
			t.Errorf("CanRevokeTarget(%q, %q) = %v, want %v", tt.actor, tt.target, got, tt.want)
		}
	}
}

func TestGetMine_UnauthorizedActor(t *testing.T) {
	svc := token.New(nil, nil)
	_, err := svc.GetMine(t.Context(), "not-a-uuid")
	if err != service.ErrUnauthorized {
		t.Errorf("GetMine() error = %v, want ErrUnauthorized", err)
	}
}

func TestCreate_UnauthorizedActor(t *testing.T) {
	svc := token.New(nil, nil)
	_, err := svc.Create(t.Context(), "")
	if err != service.ErrUnauthorized {
		t.Errorf("Create() error = %v, want ErrUnauthorized", err)
	}
}

func TestRevokeOwn_UnauthorizedActor(t *testing.T) {
	svc := token.New(nil, nil)
	err := svc.RevokeOwn(t.Context(), "bad")
	if err != service.ErrUnauthorized {
		t.Errorf("RevokeOwn() error = %v, want ErrUnauthorized", err)
	}
}

func TestRevokeForUser_ForbiddenForUserRole(t *testing.T) {
	svc := token.New(nil, nil)
	err := svc.RevokeForUser(t.Context(), "ignored", service.RoleUser, "ignored")
	if err != service.ErrForbidden {
		t.Errorf("RevokeForUser() error = %v, want ErrForbidden", err)
	}
}

func TestRevokeForUser_InvalidTargetID(t *testing.T) {
	svc := token.New(nil, nil)
	err := svc.RevokeForUser(t.Context(), "00000000-0000-0000-0000-000000000001", service.RoleAdmin, "not-a-uuid")
	if err != service.ErrUserNotFound {
		t.Errorf("RevokeForUser() error = %v, want ErrUserNotFound", err)
	}
}
