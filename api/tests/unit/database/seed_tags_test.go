package database_test

import (
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/database"
)

func TestDefaultTags_Catalog(t *testing.T) {
	want := map[string]string{
		"backend":  "#000000",
		"frontend": "#FFFFFF",
		"infra":    "#F97316",
	}
	if len(database.DefaultTags) != len(want) {
		t.Fatalf("DefaultTags len = %d, want %d", len(database.DefaultTags), len(want))
	}
	for _, tag := range database.DefaultTags {
		color, ok := want[tag.Title]
		if !ok {
			t.Errorf("unexpected default tag title %q", tag.Title)
			continue
		}
		if tag.Color != color {
			t.Errorf("tag %q color = %q, want %q", tag.Title, tag.Color, color)
		}
	}
}
