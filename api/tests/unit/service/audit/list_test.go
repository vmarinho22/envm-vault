package audit_test

import (
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
)

func TestList_ForbiddenForUserRole(t *testing.T) {
	svc := audit.New(nil, nil)
	_, err := svc.List(t.Context(), service.RoleUser, audit.ListFilter{})
	if err != service.ErrForbidden {
		t.Errorf("List() error = %v, want ErrForbidden", err)
	}
}
