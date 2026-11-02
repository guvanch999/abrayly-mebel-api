package repository

import (
	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) (*models.Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) GetAll() ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.Order("created_at, id").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) GetById(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.NewResponseError(exceptions.ErrNotFound, err)
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetByIdWithParent(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := r.db.Preload("Parent").Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Update(category *models.Category) (*models.Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Category{}).Error
}
