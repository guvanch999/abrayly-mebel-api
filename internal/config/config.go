package config

import (
	"errors"
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Schema holds every runtime setting. Values are read from the process
// environment; a local .env file (if present) is loaded first.
type Schema struct {
	Environment string `env:"GO_ENV" envDefault:"production"`
	Port        string `env:"PORT" envDefault:"3030"`
	AppUrl      string `env:"APP_URL"`

	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDb       string `env:"POSTGRES_DB"`
	PostgresSSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`

	DisableAutoMigration bool `env:"DISABLE_AUTO_MIGRATION" envDefault:"true"`

	AdminJwtSecret     string `env:"ADMIN_JWT_SECRET"`
	AdminRefreshSecret string `env:"ADMIN_REFRESH_SECRET"`

	// Seed credentials for the first admin account. The account is created
	// only when auto-migration is enabled, the admins table is empty and
	// AdminPassword is set. There is deliberately no default password.
	AdminName     string `env:"ADMIN_NAME" envDefault:"Administrator"`
	AdminUsername string `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword string `env:"ADMIN_PASSWORD"`

	// Comma-separated list of allowed CORS origins ("*" allows any origin).
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
}

// IsDevelopment reports whether the service runs in development mode.
func (s *Schema) IsDevelopment() bool {
	return s.Environment == "development"
}

// AllowedOrigins returns CorsAllowedOrigins as a trimmed list.
func (s *Schema) AllowedOrigins() []string {
	origins := make([]string, 0)
	for _, origin := range strings.Split(s.CorsAllowedOrigins, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return []string{"*"}
	}
	return origins
}

// Validate checks the settings the service cannot run safely without.
func (s *Schema) Validate() error {
	var errs []error
	if s.AdminJwtSecret == "" {
		errs = append(errs, errors.New("ADMIN_JWT_SECRET must be set"))
	}
	if s.AdminRefreshSecret == "" {
		errs = append(errs, errors.New("ADMIN_REFRESH_SECRET must be set"))
	}
	if s.AdminJwtSecret != "" && s.AdminJwtSecret == s.AdminRefreshSecret {
		errs = append(errs, errors.New("ADMIN_JWT_SECRET and ADMIN_REFRESH_SECRET must differ"))
	}
	return errors.Join(errs...)
}

// Load parses the environment (and an optional .env file) into a Schema
// without validating it.
func Load() (*Schema, error) {
	_ = godotenv.Load(".env")

	cfg := &Schema{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	cfg.applyDefaults()
	return cfg, nil
}

// applyDefaults fills settings that are present but empty (e.g. "PORT=" in
// a .env file or an empty compose variable); envDefault only covers unset
// variables.
func (s *Schema) applyDefaults() {
	setIfEmpty := func(field *string, value string) {
		if strings.TrimSpace(*field) == "" {
			*field = value
		}
	}
	setIfEmpty(&s.Environment, "production")
	setIfEmpty(&s.Port, "3030")
	setIfEmpty(&s.PostgresHost, "localhost")
	setIfEmpty(&s.PostgresPort, "5432")
	setIfEmpty(&s.PostgresSSLMode, "disable")
	setIfEmpty(&s.AdminName, "Administrator")
	setIfEmpty(&s.AdminUsername, "admin")
	setIfEmpty(&s.CorsAllowedOrigins, "*")
}

// GetConfig loads and validates the configuration, exiting on failure.
func GetConfig() *Schema {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("error parsing configuration: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}
	return cfg
}
