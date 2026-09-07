package user_test

import (
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/envm-vault/envm-vault/api/internal/service/user"
)

func TestRegister_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := user.New(nil, nil, spy)

	_, err := svc.Register(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleUser, "new@example.com", service.RoleUser)
	if err != service.ErrForbidden {
		t.Errorf("Register() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Register, got %d", len(spy.Events))
	}
}

func TestRegister_InvalidRole(t *testing.T) {
	svc := user.New(nil, nil, nil)

	_, err := svc.Register(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleSuperAdmin, "new@example.com", service.RoleSuperAdmin)
	if err != service.ErrInvalidRole {
		t.Errorf("Register() error = %v, want ErrInvalidRole", err)
	}
}

func TestList_ForbiddenForUserRole(t *testing.T) {
	svc := user.New(nil, nil, nil)

	_, err := svc.List(t.Context(), service.RoleUser)
	if err != service.ErrForbidden {
		t.Errorf("List() error = %v, want ErrForbidden", err)
	}
}

func TestBlock_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := user.New(nil, nil, spy)

	_, err := svc.Block(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000002")
	if err != service.ErrForbidden {
		t.Errorf("Block() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events, got %d", len(spy.Events))
	}
}

func TestBlock_ForbiddenSelf(t *testing.T) {
	spy := &audit.Spy{}
	svc := user.New(nil, nil, spy)
	id := "00000000-0000-0000-0000-000000000001"

	_, err := svc.Block(t.Context(), id, "admin@example.com", service.RoleAdmin, id)
	if err != service.ErrForbidden {
		t.Errorf("Block() error = %v, want ErrForbidden", err)
	}
}

func TestBlock_InvalidTargetID(t *testing.T) {
	svc := user.New(nil, nil, nil)

	_, err := svc.Block(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "not-a-uuid")
	if err != service.ErrUserNotFound {
		t.Errorf("Block() error = %v, want ErrUserNotFound", err)
	}
}

func TestUnblock_ForbiddenForUserRole(t *testing.T) {
	svc := user.New(nil, nil, nil)

	_, err := svc.Unblock(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000002")
	if err != service.ErrForbidden {
		t.Errorf("Unblock() error = %v, want ErrForbidden", err)
	}
}

func TestCanManageTarget(t *testing.T) {
	tests := []struct {
		actor  string
		target string
		want   bool
	}{
		{service.RoleSuperAdmin, service.RoleSuperAdmin, true},
		{service.RoleSuperAdmin, service.RoleAdmin, true},
		{service.RoleAdmin, service.RoleUser, true},
		{service.RoleAdmin, service.RoleAdmin, true},
		{service.RoleAdmin, service.RoleSuperAdmin, false},
		{service.RoleUser, service.RoleUser, false},
	}
	for _, tt := range tests {
		got := service.CanManageTarget(tt.actor, tt.target)
		if got != tt.want {
			t.Errorf("CanManageTarget(%q, %q) = %v, want %v", tt.actor, tt.target, got, tt.want)
		}
	}
}
