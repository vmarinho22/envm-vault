package secret_test

import (
	"strings"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/secretcrypto"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/secret"
	"github.com/google/uuid"
)

func testCrypter(t *testing.T) secretcrypto.Crypter {
	t.Helper()
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 3)
	}
	c, err := secretcrypto.NewAESGCM(key, 1)
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	return c
}

func TestCreate_InvalidProjectID(t *testing.T) {
	svc := secret.New(nil, nil, testCrypter(t), nil)
	_, err := svc.Create(t.Context(), uuid.Nil.String(), "a@b.com", service.RoleUser, "not-a-uuid", "OK_NAME", nil)
	if err != service.ErrProjectNotFound {
		t.Errorf("Create() error = %v, want ErrProjectNotFound", err)
	}
}

func TestUpdateName_InvalidProjectID(t *testing.T) {
	svc := secret.New(nil, nil, testCrypter(t), nil)
	_, err := svc.UpdateName(t.Context(), uuid.Nil.String(), "a@b.com", service.RoleUser, "bad", "bad", "NAME")
	if err != service.ErrProjectNotFound {
		t.Errorf("UpdateName() error = %v, want ErrProjectNotFound", err)
	}
}

func TestDelete_InvalidProjectID(t *testing.T) {
	svc := secret.New(nil, nil, testCrypter(t), nil)
	err := svc.Delete(t.Context(), uuid.Nil.String(), "a@b.com", service.RoleUser, "bad", "bad")
	if err != service.ErrProjectNotFound {
		t.Errorf("Delete() error = %v, want ErrProjectNotFound", err)
	}
}

func TestFormatDotEnv(t *testing.T) {
	out := secret.FormatDotEnv([]secret.DownloadItem{
		{Name: "A", Value: "simple"},
		{Name: "B", Value: "has space"},
		{Name: "C", Value: `quote"here`},
		{Name: "D", Value: ""},
	})
	want := "A=simple\nB=\"has space\"\nC=\"quote\\\"here\"\nD=\"\"\n"
	if out != want {
		t.Fatalf("FormatDotEnv = %q, want %q", out, want)
	}
}

func TestFormatDotEnv_LongNameRejectedByPatternElsewhere(t *testing.T) {
	// Document max length constant via FormatDotEnv still works with long names.
	name := strings.Repeat("A", 128)
	out := secret.FormatDotEnv([]secret.DownloadItem{{Name: name, Value: "v"}})
	if !strings.HasPrefix(out, name+"=") {
		t.Fatalf("unexpected output prefix")
	}
}

func TestDownloadBySlug_RequiresEnvironment(t *testing.T) {
	svc := secret.New(nil, nil, testCrypter(t), nil)
	_, err := svc.DownloadBySlug(t.Context(), uuid.Nil.String(), "a@b.com", service.RoleUser, "my-app", "")
	if err != service.ErrEnvironmentParamRequired {
		t.Errorf("DownloadBySlug() error = %v, want ErrEnvironmentParamRequired", err)
	}
}

func TestDownloadBySlug_ProductionForbiddenForNonSuperAdmin(t *testing.T) {
	svc := secret.New(nil, nil, testCrypter(t), nil)
	for _, role := range []string{service.RoleAdmin, service.RoleUser} {
		_, err := svc.DownloadBySlug(t.Context(), uuid.Nil.String(), "a@b.com", role, "my-app", "production")
		if err != service.ErrProductionDownloadForbidden {
			t.Errorf("role %s: error = %v, want ErrProductionDownloadForbidden", role, err)
		}
	}
}
