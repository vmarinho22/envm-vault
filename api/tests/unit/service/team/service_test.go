package team_test

import (
	"testing"

	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/vmarinho22/envm-vault/api/internal/service/audit"
	"github.com/vmarinho22/envm-vault/api/internal/service/team"
)

func TestCreate_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := team.New(nil, nil, spy)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "Platform", "00000000-0000-0000-0000-000000000002")
	if err != service.ErrForbidden {
		t.Errorf("Create() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Create, got %d", len(spy.Events))
	}
}

func TestUpdate_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := team.New(nil, nil, spy)
	name := "Platform"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000099", &name, nil)
	if err != service.ErrForbidden {
		t.Errorf("Update() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Update, got %d", len(spy.Events))
	}
}

func TestDelete_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := team.New(nil, nil, spy)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000099")
	if err != service.ErrForbidden {
		t.Errorf("Delete() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Delete, got %d", len(spy.Events))
	}
}

func TestCreate_InvalidTeamName(t *testing.T) {
	svc := team.New(nil, nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "!!!", "00000000-0000-0000-0000-000000000002")
	if err != service.ErrInvalidTeamName {
		t.Errorf("Create() error = %v, want ErrInvalidTeamName", err)
	}
}

func TestCreate_InvalidAccountableID(t *testing.T) {
	svc := team.New(nil, nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "Platform", "not-a-uuid")
	if err != service.ErrUserNotFound {
		t.Errorf("Create() error = %v, want ErrUserNotFound", err)
	}
}

func TestUpdate_InvalidTeamID(t *testing.T) {
	svc := team.New(nil, nil, nil)
	name := "Platform"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id", &name, nil)
	if err != service.ErrTeamNotFound {
		t.Errorf("Update() error = %v, want ErrTeamNotFound", err)
	}
}

func TestDelete_InvalidTeamID(t *testing.T) {
	svc := team.New(nil, nil, nil)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id")
	if err != service.ErrTeamNotFound {
		t.Errorf("Delete() error = %v, want ErrTeamNotFound", err)
	}
}
