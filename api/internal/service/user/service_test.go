package user

import (
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service"
)

func TestRegister_ForbiddenForUserRole(t *testing.T) {
	svc := &Service{}

	_, err := svc.Register(t.Context(), service.RoleUser, "new@example.com", service.RoleUser)
	if err != service.ErrForbidden {
		t.Errorf("Register() error = %v, want ErrForbidden", err)
	}
}

func TestRegister_InvalidRole(t *testing.T) {
	svc := &Service{}

	_, err := svc.Register(t.Context(), service.RoleSuperAdmin, "new@example.com", service.RoleSuperAdmin)
	if err != service.ErrInvalidRole {
		t.Errorf("Register() error = %v, want ErrInvalidRole", err)
	}
}
