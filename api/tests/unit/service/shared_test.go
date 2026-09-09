package service_test

import (
	"testing"

	"github.com/vmarinho22/envm-vault/api/internal/service"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid email lowercase",
			input: "user@example.com",
			want:  "user@example.com",
		},
		{
			name:  "valid email uppercase normalized",
			input: "User@Example.COM",
			want:  "user@example.com",
		},
		{
			name:  "valid email with spaces trimmed",
			input: "  user@example.com  ",
			want:  "user@example.com",
		},
		{
			name:    "invalid email",
			input:   "not-an-email",
			wantErr: true,
		},
		{
			name:    "empty email",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.NormalizeEmail(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("NormalizeEmail() expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NormalizeEmail() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NormalizeEmail() = %q, want %q", got, tt.want)
			}
		})
	}
}
