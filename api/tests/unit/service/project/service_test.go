package project_test

import (
	"strings"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/envm-vault/envm-vault/api/internal/service/project"
)

func TestDelete_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := project.New(nil, nil, nil, nil, spy)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000099")
	if err != service.ErrForbidden {
		t.Errorf("Delete() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Delete, got %d", len(spy.Events))
	}
}

func TestCreate_AllowsUserRole_InvalidTitle(t *testing.T) {
	spy := &audit.Spy{}
	svc := project.New(nil, nil, nil, nil, spy)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "!!!", []string{"00000000-0000-0000-0000-000000000002"}, []string{"00000000-0000-0000-0000-000000000003"}, nil)
	if err != service.ErrInvalidProjectTitle {
		t.Errorf("Create() error = %v, want ErrInvalidProjectTitle", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events, got %d", len(spy.Events))
	}
}

func TestCreate_TitleTooLong(t *testing.T) {
	svc := project.New(nil, nil, nil, nil, nil)
	long := strings.Repeat("a", 81)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, long, []string{"00000000-0000-0000-0000-000000000002"}, []string{"00000000-0000-0000-0000-000000000003"}, nil)
	if err != service.ErrInvalidProjectTitle {
		t.Errorf("Create() error = %v, want ErrInvalidProjectTitle", err)
	}
}

func TestCreate_MissingRelations(t *testing.T) {
	svc := project.New(nil, nil, nil, nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "My App", nil, []string{"00000000-0000-0000-0000-000000000003"}, nil)
	if err != service.ErrInvalidProjectRelations {
		t.Errorf("Create() error = %v, want ErrInvalidProjectRelations", err)
	}
}

func TestUpdate_InvalidProjectID(t *testing.T) {
	svc := project.New(nil, nil, nil, nil, nil)
	title := "My App"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "bad-id", &title, nil, nil, nil)
	if err != service.ErrProjectNotFound {
		t.Errorf("Update() error = %v, want ErrProjectNotFound", err)
	}
}

func TestDelete_InvalidProjectID(t *testing.T) {
	svc := project.New(nil, nil, nil, nil, nil)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id")
	if err != service.ErrProjectNotFound {
		t.Errorf("Delete() error = %v, want ErrProjectNotFound", err)
	}
}

func TestDefaultPageSize(t *testing.T) {
	if project.DefaultPageSize != 12 {
		t.Errorf("DefaultPageSize = %d, want 12", project.DefaultPageSize)
	}
}
