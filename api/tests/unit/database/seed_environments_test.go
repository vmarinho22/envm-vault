package database_test

import (
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/database"
)

func TestDefaultEnvironments_Catalog(t *testing.T) {
	want := []string{"dev", "test", "staging", "production"}
	if len(database.DefaultEnvironments) != len(want) {
		t.Fatalf("DefaultEnvironments len = %d, want %d", len(database.DefaultEnvironments), len(want))
	}
	for i, title := range want {
		got := database.DefaultEnvironments[i].Title
		if got != title {
			t.Errorf("DefaultEnvironments[%d].Title = %q, want %q", i, got, title)
		}
	}
}
