package slug_test

import (
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/slug"
)

func TestFromName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "simple", in: "Platform", want: "platform"},
		{name: "accents", in: "Equipe São Paulo", want: "equipe-sao-paulo"},
		{name: "punctuation", in: "Foo & Bar!!!", want: "foo-bar"},
		{name: "spaces and hyphens", in: "  My -- Team  ", want: "my-team"},
		{name: "numbers", in: "Team 42", want: "team-42"},
		{name: "empty", in: "   ", want: ""},
		{name: "only symbols", in: "!!!", want: ""},
		{name: "cedilla", in: "Ação", want: "acao"},
		{name: "umlaut", in: "Über Team", want: "uber-team"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slug.FromName(tt.in)
			if got != tt.want {
				t.Errorf("FromName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
