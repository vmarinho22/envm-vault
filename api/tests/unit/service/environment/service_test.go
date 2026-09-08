package environment_test

import (
	"strings"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/envm-vault/envm-vault/api/internal/service/environment"
)

func TestCreate_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := environment.New(nil, spy)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "qa")
	if err != service.ErrForbidden {
		t.Errorf("Create() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Create, got %d", len(spy.Events))
	}
}

func TestUpdate_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := environment.New(nil, spy)
	title := "qa"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000099", &title)
	if err != service.ErrForbidden {
		t.Errorf("Update() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Update, got %d", len(spy.Events))
	}
}

func TestDelete_ForbiddenForUserRole(t *testing.T) {
	spy := &audit.Spy{}
	svc := environment.New(nil, spy)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "00000000-0000-0000-0000-000000000099")
	if err != service.ErrForbidden {
		t.Errorf("Delete() error = %v, want ErrForbidden", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on forbidden Delete, got %d", len(spy.Events))
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	svc := environment.New(nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "   ")
	if err != service.ErrInvalidEnvironmentTitle {
		t.Errorf("Create() error = %v, want ErrInvalidEnvironmentTitle", err)
	}
}

func TestCreate_TitleTooLong(t *testing.T) {
	svc := environment.New(nil, nil)
	long := strings.Repeat("a", 81)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, long)
	if err != service.ErrInvalidEnvironmentTitle {
		t.Errorf("Create() error = %v, want ErrInvalidEnvironmentTitle", err)
	}
}

func TestUpdate_InvalidEnvironmentID(t *testing.T) {
	svc := environment.New(nil, nil)
	title := "qa"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id", &title)
	if err != service.ErrEnvironmentNotFound {
		t.Errorf("Update() error = %v, want ErrEnvironmentNotFound", err)
	}
}

func TestDelete_InvalidEnvironmentID(t *testing.T) {
	svc := environment.New(nil, nil)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id")
	if err != service.ErrEnvironmentNotFound {
		t.Errorf("Delete() error = %v, want ErrEnvironmentNotFound", err)
	}
}

func TestAssertMutableEnvironment_ProductionProtected(t *testing.T) {
	cases := []struct {
		title string
		want  error
	}{
		{title: service.ProductionEnvironmentTitle, want: service.ErrProductionEnvironmentProtected},
		{title: "dev", want: nil},
		{title: "staging", want: nil},
	}
	for _, tc := range cases {
		err := environment.AssertMutableEnvironment(tc.title)
		if err != tc.want {
			t.Errorf("AssertMutableEnvironment(%q) = %v, want %v", tc.title, err, tc.want)
		}
	}
}
