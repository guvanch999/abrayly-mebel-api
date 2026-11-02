package service

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
)

const (
	categoryImagesDir = "uploads/categories"
	// maxCategoryDepth bounds the ancestor walk when re-parenting a category.
	maxCategoryDepth = 100
)

type CategoryService struct {
	repository *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repository: repo,
	}
}

func (s *CategoryService) CreateCategory(categoryDto *dto.CategoryDto, file *multipart.FileHeader) (*models.Category, error) {
	category := &models.Category{
		ID:     uuid.New(),
		NameTm: categoryDto.NameTm,
		NameRu: categoryDto.NameRu,
		NameEn: categoryDto.NameEn,
	}

	if categoryDto.ParentId != nil {
		if _, err := s.repository.GetById(*categoryDto.ParentId); err != nil {
			return nil, err
		}
		category.ParentId = categoryDto.ParentId
	}

	filePath, err := util.SaveUploadedFile(file, category.ID.String(), categoryImagesDir)
	if err != nil {
		return nil, err
	}
	category.ImageUrl = filePath

	created, err := s.repository.Create(category)
	if err != nil {
		_ = util.DeleteFile(filePath)
		return nil, err
	}
	return created, nil
}

func (s *CategoryService) GetCategory(id uuid.UUID) (*models.Category, error) {
	return s.repository.GetById(id)
}

// GetAllCategories returns every category arranged as a tree.
func (s *CategoryService) GetAllCategories() ([]*dto.CategoryResponseDto, error) {
	categories, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return dto.BuildCategoryTree(categories), nil
}

func (s *CategoryService) UpdateCategory(id uuid.UUID, categoryDto *dto.CategoryDto, file *multipart.FileHeader) (*models.Category, error) {
	category, err := s.repository.GetById(id)
	if err != nil {
		return nil, err
	}

	if categoryDto.ParentId != nil {
		if err := s.checkParent(id, *categoryDto.ParentId); err != nil {
			return nil, err
		}
		category.ParentId = categoryDto.ParentId
	}

	if file != nil {
		filePath, err := util.SaveUploadedFile(file, category.ID.String(), categoryImagesDir)
		if err != nil {
			return nil, err
		}
		if category.ImageUrl != "" && category.ImageUrl != filePath {
			_ = util.DeleteFile(category.ImageUrl)
		}
		category.ImageUrl = filePath
	}

	category.NameTm = categoryDto.NameTm
	category.NameRu = categoryDto.NameRu
	category.NameEn = categoryDto.NameEn

	return s.repository.Update(category)
}

// checkParent makes sure parentID exists and is neither the category itself
// nor one of its descendants, which would detach a cycle from the tree.
func (s *CategoryService) checkParent(id, parentID uuid.UUID) error {
	errCycle := exceptions.NewResponseError(exceptions.ErrBadRequest, errors.New("a category cannot be moved under itself or its subcategory"))

	current := &parentID
	for depth := 0; current != nil; depth++ {
		if *current == id || depth > maxCategoryDepth {
			return errCycle
		}
		ancestor, err := s.repository.GetById(*current)
		if err != nil {
			return err
		}
		current = ancestor.ParentId
	}
	return nil
}

// DeleteCategory deletes a category; subcategories are removed by the
// ON DELETE CASCADE constraint and products keep existing with category_id
// set to NULL.
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	if _, err := s.repository.GetById(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}
