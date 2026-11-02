package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/guvanch999/abrayly-mebel-backend/internal/config"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// DSN builds a PostgreSQL connection URL from the configuration.
// Credentials are URL-escaped so special characters in passwords are safe.
func DSN(cfg *config.Schema) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.PostgresUser, cfg.PostgresPassword),
		Host:     fmt.Sprintf("%s:%s", cfg.PostgresHost, cfg.PostgresPort),
		Path:     cfg.PostgresDb,
		RawQuery: url.Values{"sslmode": []string{cfg.PostgresSSLMode}}.Encode(),
	}
	return u.String()
}

// Connect opens the connection pool, verifies connectivity and, unless
// DISABLE_AUTO_MIGRATION is set, migrates the schema and seeds the first admin.
func Connect(cfg *config.Schema) (*gorm.DB, error) {
	logLevel := gormLogger.Warn
	if cfg.IsDevelopment() {
		logLevel = gormLogger.Info
	}

	database, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  DSN(cfg),
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger:          gormLogger.Default.LogMode(logLevel),
			NowFunc:         func() time.Time { return time.Now().UTC() },
			PrepareStmt:     false,
			CreateBatchSize: 100,
		},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(40)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if !cfg.DisableAutoMigration {
		if err := migrate(database); err != nil {
			return nil, err
		}

		seed := AdminSeed{Name: cfg.AdminName, Username: cfg.AdminUsername, Password: cfg.AdminPassword}
		if err := SeedAdmin(database, seed); err != nil {
			return nil, fmt.Errorf("error seeding admin: %w", err)
		}
	}

	log.Println("Connected to PostgreSQL")
	return database, nil
}

func migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return fmt.Errorf("error creating extension: %w", err)
	}
	err := db.AutoMigrate(
		&models.Category{},
		&models.Admin{},
		&models.Product{},
		&models.ProductImage{},
		&models.Attribute{},
		&models.Banner{},
	)
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}
	return nil
}
