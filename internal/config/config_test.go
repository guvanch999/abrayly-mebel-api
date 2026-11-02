package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Schema
		wantErr []string
	}{
		{
			name: "valid",
			cfg:  Schema{AdminJwtSecret: "access", AdminRefreshSecret: "refresh"},
		},
		{
			name:    "missing secrets",
			cfg:     Schema{},
			wantErr: []string{"ADMIN_JWT_SECRET must be set", "ADMIN_REFRESH_SECRET must be set"},
		},
		{
			name:    "missing refresh secret",
			cfg:     Schema{AdminJwtSecret: "access"},
			wantErr: []string{"ADMIN_REFRESH_SECRET must be set"},
		},
		{
			name:    "identical secrets",
			cfg:     Schema{AdminJwtSecret: "same", AdminRefreshSecret: "same"},
			wantErr: []string{"must differ"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if len(tt.wantErr) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected an error")
			}
			for _, want := range tt.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not mention %q", err, want)
				}
			}
		})
	}
}

func TestAllowedOrigins(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{in: "*", want: []string{"*"}},
		{in: "", want: []string{"*"}},
		{in: " , ", want: []string{"*"}},
		{in: "https://abrayly.com", want: []string{"https://abrayly.com"}},
		{in: "https://a.example, https://b.example ,", want: []string{"https://a.example", "https://b.example"}},
	}
	for _, tt := range tests {
		cfg := Schema{CorsAllowedOrigins: tt.in}
		if got := cfg.AllowedOrigins(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("AllowedOrigins(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

// chdirTemp runs the test from an empty directory so a developer's .env
// file cannot leak into the result.
func chdirTemp(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
}

// unsetenv removes key for the duration of the test.
func unsetenv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "") // registers restoration of the original value
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv(%s): %v", key, err)
	}
}

var allKeys = []string{
	"GO_ENV", "PORT", "APP_URL",
	"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_DB", "POSTGRES_SSLMODE",
	"DISABLE_AUTO_MIGRATION", "ADMIN_JWT_SECRET", "ADMIN_REFRESH_SECRET",
	"ADMIN_NAME", "ADMIN_USERNAME", "ADMIN_PASSWORD", "CORS_ALLOWED_ORIGINS",
}

func TestLoadDefaults(t *testing.T) {
	chdirTemp(t)
	for _, key := range allKeys {
		unsetenv(t, key)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := Schema{
		Environment:          "production",
		Port:                 "3030",
		PostgresHost:         "localhost",
		PostgresPort:         "5432",
		PostgresSSLMode:      "disable",
		DisableAutoMigration: true,
		AdminName:            "Administrator",
		AdminUsername:        "admin",
		CorsAllowedOrigins:   "*",
	}
	if !reflect.DeepEqual(*cfg, want) {
		t.Errorf("Load() = %+v, want %+v", *cfg, want)
	}
	if cfg.AdminPassword != "" {
		t.Error("there must be no default admin password")
	}
}

func TestLoadReadsEnvironment(t *testing.T) {
	chdirTemp(t)
	for _, key := range allKeys {
		unsetenv(t, key)
	}
	t.Setenv("GO_ENV", "development")
	t.Setenv("PORT", "8080")
	t.Setenv("POSTGRES_DB", "abrayly")
	t.Setenv("DISABLE_AUTO_MIGRATION", "false")
	t.Setenv("ADMIN_USERNAME", "owner")
	t.Setenv("ADMIN_PASSWORD", "from-env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !cfg.IsDevelopment() || cfg.Port != "8080" || cfg.PostgresDb != "abrayly" {
		t.Errorf("unexpected config: env=%q port=%q db=%q", cfg.Environment, cfg.Port, cfg.PostgresDb)
	}
	if cfg.DisableAutoMigration {
		t.Error("DisableAutoMigration = true, want false")
	}
	if cfg.AdminUsername != "owner" || cfg.AdminPassword != "from-env" {
		t.Error("admin seed credentials were not read from the environment")
	}
}

func TestLoadTreatsEmptyValuesAsDefaults(t *testing.T) {
	chdirTemp(t)
	for _, key := range allKeys {
		unsetenv(t, key)
	}
	t.Setenv("PORT", "")
	t.Setenv("POSTGRES_HOST", " ")
	t.Setenv("ADMIN_USERNAME", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != "3030" || cfg.PostgresHost != "localhost" || cfg.AdminUsername != "admin" {
		t.Errorf("empty values were not replaced: port=%q host=%q user=%q", cfg.Port, cfg.PostgresHost, cfg.AdminUsername)
	}
}

func TestIsDevelopment(t *testing.T) {
	if !(&Schema{Environment: "development"}).IsDevelopment() {
		t.Error("development should be reported as development")
	}
	if (&Schema{Environment: "production"}).IsDevelopment() {
		t.Error("production should not be reported as development")
	}
}
