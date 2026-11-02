package dto

import (
	"testing"

	"github.com/guvanch999/abrayly-mebel-backend/utils/validator"
)

func TestLoginRequestValidation(t *testing.T) {
	v := validator.NewValidator()

	tests := []struct {
		name    string
		in      LoginRequestDto
		wantErr bool
	}{
		{name: "valid", in: LoginRequestDto{Username: "admin", Password: "secret"}},
		{name: "missing username", in: LoginRequestDto{Password: "secret"}, wantErr: true},
		{name: "missing password", in: LoginRequestDto{Username: "admin"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(&tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRefreshTokenRequestValidation(t *testing.T) {
	v := validator.NewValidator()

	if err := v.Validate(&RefreshTokenRequestDto{RefreshToken: "a.b.c"}); err != nil {
		t.Errorf("valid request rejected: %v", err)
	}
	if err := v.Validate(&RefreshTokenRequestDto{}); err == nil {
		t.Error("empty refresh token accepted")
	}
}
