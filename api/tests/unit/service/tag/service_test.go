package tag_test

import (
	"strings"
	"testing"

	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/vmarinho22/envm-vault/api/internal/service/audit"
	"github.com/vmarinho22/envm-vault/api/internal/service/tag"
)

func TestCreate_AllowsUserRole_InvalidTitle(t *testing.T) {
	spy := &audit.Spy{}
	svc := tag.New(nil, spy)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "", "#6366F1")
	if err != service.ErrInvalidTagTitle {
		t.Errorf("Create() error = %v, want ErrInvalidTagTitle (user role must reach validation)", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on invalid Create, got %d", len(spy.Events))
	}
}

func TestUpdate_AllowsUserRole_InvalidTagID(t *testing.T) {
	spy := &audit.Spy{}
	svc := tag.New(nil, spy)
	title := "Backend"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "bad-id", &title, nil)
	if err != service.ErrTagNotFound {
		t.Errorf("Update() error = %v, want ErrTagNotFound (user role must reach lookup)", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on invalid Update, got %d", len(spy.Events))
	}
}

func TestDelete_AllowsUserRole_InvalidTagID(t *testing.T) {
	spy := &audit.Spy{}
	svc := tag.New(nil, spy)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "bad-id")
	if err != service.ErrTagNotFound {
		t.Errorf("Delete() error = %v, want ErrTagNotFound (user role must reach lookup)", err)
	}
	if len(spy.Events) != 0 {
		t.Errorf("expected no audit events on invalid Delete, got %d", len(spy.Events))
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	svc := tag.New(nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "   ", "#6366F1")
	if err != service.ErrInvalidTagTitle {
		t.Errorf("Create() error = %v, want ErrInvalidTagTitle", err)
	}
}

func TestCreate_TitleTooLong(t *testing.T) {
	svc := tag.New(nil, nil)
	long := strings.Repeat("a", 81)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, long, "#6366F1")
	if err != service.ErrInvalidTagTitle {
		t.Errorf("Create() error = %v, want ErrInvalidTagTitle", err)
	}
}

func TestCreate_InvalidColor(t *testing.T) {
	svc := tag.New(nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "Backend", "red")
	if err != service.ErrInvalidTagColor {
		t.Errorf("Create() error = %v, want ErrInvalidTagColor", err)
	}
}

func TestCreate_InvalidColorShortHex(t *testing.T) {
	svc := tag.New(nil, nil)

	_, err := svc.Create(t.Context(), "00000000-0000-0000-0000-000000000001", "user@example.com", service.RoleUser, "Backend", "#FFF")
	if err != service.ErrInvalidTagColor {
		t.Errorf("Create() error = %v, want ErrInvalidTagColor", err)
	}
}

func TestUpdate_InvalidTagID(t *testing.T) {
	svc := tag.New(nil, nil)
	title := "Backend"

	_, err := svc.Update(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id", &title, nil)
	if err != service.ErrTagNotFound {
		t.Errorf("Update() error = %v, want ErrTagNotFound", err)
	}
}

func TestDelete_InvalidTagID(t *testing.T) {
	svc := tag.New(nil, nil)

	err := svc.Delete(t.Context(), "00000000-0000-0000-0000-000000000001", "admin@example.com", service.RoleAdmin, "bad-id")
	if err != service.ErrTagNotFound {
		t.Errorf("Delete() error = %v, want ErrTagNotFound", err)
	}
}
