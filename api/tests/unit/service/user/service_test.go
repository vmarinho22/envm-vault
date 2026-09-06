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
