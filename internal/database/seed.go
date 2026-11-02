package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"gorm.io/gorm"
)

// AdminSeed holds the credentials used to create the first admin account.
type AdminSeed struct {
	Name     string
	Username string
	Password string
}

// SeedAdmin creates the first admin account when the admins table is empty.
//
// Credentials come exclusively from the environment (ADMIN_USERNAME /
// ADMIN_PASSWORD). If no password is configured, seeding is skipped with a
// warning: the service never falls back to a built-in default password.
func SeedAdmin(db *gorm.DB, seed AdminSeed) error {
	if seed.Password == "" {
		log.Println("WARNING: ADMIN_PASSWORD is not set; skipping default admin seeding")
		return nil
	}

	username := strings.TrimSpace(seed.Username)
	if username == "" {
		username = "admin"
	}
	name := strings.TrimSpace(seed.Name)
	if name == "" {
		name = "Administrator"
	}

	var count int64
	if err := db.Model(&models.Admin{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count admins: %w", err)
	}
	if count > 0 {
		return nil
	}

	hashedPassword, err := util.HashString(seed.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := models.Admin{
		ID:       uuid.New(),
		Name:     name,
		Username: username,
		Password: hashedPassword,
	}
	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	log.Printf("Seeded initial admin account %q", username)
	return nil
}
