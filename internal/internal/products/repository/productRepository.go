package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// catalogOrder lists newest products first; id breaks ties so pages are stable.
const catalogOrder = "created_at DESC, id"

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FilterScope applies the catalog filters:
//   - category_id: exact category match;
//   - search: case-insensitive substring match on the Turkmen, Russian or
//     English name (ILIKE, with LIKE wildcards in the input escaped);
//   - is_discount=true: only products with a discount.
func FilterScope(filter *dto.FilterDto) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		if filter.CategoryID != nil {
			db = db.Where("category_id = ?", *filter.CategoryID)
		}
		if term := filter.SearchTerm(); term != "" {
			pattern := "%" + EscapeLike(term) + "%"
			// GORM wraps an OR expression in parentheses when it is combined
			// with other conditions, so the AND filters still apply.
			db = db.Where("name_tm ILIKE ? OR name_ru ILIKE ? OR name_en ILIKE ?", pattern, pattern, pattern)
		}
		if filter.OnlyDiscounted() {
			db = db.Where("discount > 0")
		}
		return db
	}
}

// PaginateScope applies LIMIT/OFFSET for a normalized PaginateDto.
func PaginateScope(paginate *dto.PaginateDto) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(paginate.Offset()).Limit(paginate.PageSize)
	}
}

// EscapeLike escapes the LIKE/ILIKE wildcards % and _ (and the escape
// character itself) so user input is matched literally.
func EscapeLike(s string) string {
	return likeEscaper.Replace(s)
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (r *ProductRepository) CreateProduct(product *models.Product) (*models.Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) GetProduct(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).Preload("Attributes").Preload("Category").Preload("Images").First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewResponseError(exceptions.ErrNotFound, errors.New("product not found"))
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetCategory(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewResponseError(exceptions.ErrNotFound, errors.New("category not found: "+id.String()))
		}
		return nil, err
	}
	return &category, nil
}

// UpdateProduct saves the product's own columns; associations (images,
// attributes, category) are managed through their dedicated endpoints.
func (r *ProductRepository) UpdateProduct(product *models.Product) (*models.Product, error) {
	if err := r.db.Omit(clause.Associations).Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// DeleteProduct removes the product together with its images and
// attributes in a single transaction.
func (r *ProductRepository) DeleteProduct(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.ProductImage{}, "product_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Attribute{}, "product_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Product{}, "id = ?", id).Error
	})
}

// GetProducts returns one page of products matching filter plus the total
// number of matches.
func (r *ProductRepository) GetProducts(paginate *dto.PaginateDto, filter *dto.FilterDto) ([]models.Product, int, error) {
	var total int64
	err := r.db.Model(&models.Product{}).Scopes(FilterScope(filter)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	products := make([]models.Product, 0)
	if total == 0 {
		return products, 0, nil
	}

	err = r.db.Model(&models.Product{}).
		Scopes(FilterScope(filter), PaginateScope(paginate)).
		Order(catalogOrder).
		Preload("Attributes").
		Preload("Images").
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

// GetRecommendedProducts returns other products from the same category.
func (r *ProductRepository) GetRecommendedProducts(paginate *dto.PaginateDto, product *models.Product) ([]models.Product, int, error) {
	sameCategory := func(db *gorm.DB) *gorm.DB {
		return db.Where("category_id = ? AND id <> ?", product.CategoryID, product.ID)
	}

	var total int64
	if err := r.db.Model(&models.Product{}).Scopes(sameCategory).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	products := make([]models.Product, 0)
	if total == 0 {
		return products, 0, nil
	}

	err := r.db.Model(&models.Product{}).
		Scopes(sameCategory, PaginateScope(paginate)).
		Order(catalogOrder).
		Preload("Attributes").
		Preload("Images").
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, int(total), nil
}

func (r *ProductRepository) AddImageToProduct(image *models.ProductImage) error {
	return r.db.Create(image).Error
}

// GetImage returns an image that belongs to the given product.
func (r *ProductRepository) GetImage(productID, imageID uuid.UUID) (*models.ProductImage, error) {
	var image models.ProductImage
	err := r.db.Where("id = ? AND product_id = ?", imageID, productID).First(&image).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewResponseError(exceptions.ErrNotFound, errors.New("image not found"))
		}
		return nil, err
	}
	return &image, nil
}

func (r *ProductRepository) RemoveImageFromProduct(productID, imageID uuid.UUID) error {
	return r.db.Delete(&models.ProductImage{}, "id = ? AND product_id = ?", imageID, productID).Error
}

func (r *ProductRepository) AddAttributeToProduct(attribute *models.Attribute) error {
	return r.db.Create(attribute).Error
}

func (r *ProductRepository) RemoveAttributeFromProduct(productID, attributeID uuid.UUID) error {
	res := r.db.Delete(&models.Attribute{}, "id = ? AND product_id = ?", attributeID, productID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return exceptions.NewResponseError(exceptions.ErrNotFound, errors.New("attribute not found"))
	}
	return nil
}

// GetAttribute returns an attribute that belongs to the given product.
func (r *ProductRepository) GetAttribute(productID, attributeID uuid.UUID) (*models.Attribute, error) {
	var attribute models.Attribute
	err := r.db.Where("id = ? AND product_id = ?", attributeID, productID).First(&attribute).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewResponseError(exceptions.ErrNotFound, errors.New("attribute not found"))
		}
		return nil, err
	}
	return &attribute, nil
}

// SaveAttribute persists an attribute, including its JSON-serialized values.
func (r *ProductRepository) SaveAttribute(attribute *models.Attribute) error {
	return r.db.Save(attribute).Error
}
