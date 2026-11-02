package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"gorm.io/gorm"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) Create(banner *models.Banner) (*models.Banner, error) {
	err := r.db.Create(banner).Error
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (r *BannerRepository) GetAll() ([]*models.Banner, error) {
	var banners []*models.Banner
	err := r.db.Find(&banners).Error
	if err != nil {
		return nil, err
	}
	return banners, nil
}

func (r *BannerRepository) Update(banner *models.Banner) (*models.Banner, error) {
	err := r.db.Save(banner).Error
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (r *BannerRepository) GetByID(id uuid.UUID) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.Where("id = ?", id).First(&banner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewResponseError(exceptions.ErrNotFound, err)
		}
		return nil, err
	}
	return &banner, nil
}

func (r *BannerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Banner{}, id).Error
}
