package security_test

import (
	"testing"

	"github.com/roushou/pocpoc/internal/security"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "simple password",
			password: "password",
			wantErr:  false,
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!2023",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, hashErr := security.HashPassword(tt.password)
			hasHashErr := hashErr != nil
			if hasHashErr != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", hashErr, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("HashPassword() returned an empty string when no error was expected")
			}
		})
	}
}
