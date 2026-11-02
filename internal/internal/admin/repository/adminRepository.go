package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"gorm.io/gorm"
)

// ErrAdminNotFound is returned when no admin matches the lookup.
var ErrAdminNotFound = errors.New("admin not found")

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, notFound(err)
	}
	return &admin, nil
}

func (r *AdminRepository) FindByID(id uuid.UUID) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, notFound(err)
	}
	return &admin, nil
}

func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAdminNotFound
	}
	return err
}
