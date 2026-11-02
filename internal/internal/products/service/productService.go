package service

import (
	"mime/multipart"
	"slices"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
)

const productImagesDir = "uploads/products"

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) CreateProduct(product *dto.ProductDto) (*models.Product, error) {
	if _, err := s.productRepository.GetCategory(product.CategoryID); err != nil {
		return nil, err
	}

	productModel := &models.Product{ID: uuid.New()}
	applyProductDto(productModel, product)

	return s.productRepository.CreateProduct(productModel)
}

func (s *ProductService) UpdateProduct(productID uuid.UUID, product *dto.ProductDto) (*models.Product, error) {
	existing, err := s.productRepository.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	if product.CategoryID != existing.CategoryID {
		category, err := s.productRepository.GetCategory(product.CategoryID)
		if err != nil {
			return nil, err
		}
		existing.Category = category
	}

	applyProductDto(existing, product)

	return s.productRepository.UpdateProduct(existing)
}

func applyProductDto(product *models.Product, input *dto.ProductDto) {
	product.CategoryID = input.CategoryID
	product.NameTm = input.NameTm
	product.NameRu = input.NameRu
	product.NameEn = input.NameEn
	product.Price = input.Price
	product.Discount = input.Discount
	product.OldPrice = input.OldPrice
}

func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	product, err := s.productRepository.GetProduct(id)
	if err != nil {
		return err
	}

	if err := s.productRepository.DeleteProduct(id); err != nil {
		return err
	}

	for _, image := range product.Images {
		_ = util.DeleteFile(image.ImageURL)
	}
	return nil
}

func (s *ProductService) GetProduct(id uuid.UUID) (*models.Product, error) {
	return s.productRepository.GetProduct(id)
}

func (s *ProductService) GetProducts(paginate *dto.PaginateDto, filter *dto.FilterDto) (*dto.ProductResponseDto, error) {
	paginate.Normalize()

	categoryName := ""
	if filter.CategoryID != nil {
		category, err := s.productRepository.GetCategory(*filter.CategoryID)
		if err != nil {
			return nil, err
		}
		categoryName = category.NameTm
	}

	products, total, err := s.productRepository.GetProducts(paginate, filter)
	if err != nil {
		return nil, err
	}

	return &dto.ProductResponseDto{
		Total:        total,
		CategoryName: categoryName,
		Data:         products,
	}, nil
}

func (s *ProductService) AddImageToProduct(productID uuid.UUID, image *multipart.FileHeader) error {
	if _, err := s.productRepository.GetProduct(productID); err != nil {
		return err
	}

	imageID := uuid.New()
	imageURL, err := util.SaveUploadedFile(image, imageID.String(), productImagesDir)
	if err != nil {
		return err
	}

	imageModel := &models.ProductImage{
		ID:        imageID,
		ProductID: productID,
		ImageURL:  imageURL,
	}
	if err := s.productRepository.AddImageToProduct(imageModel); err != nil {
		_ = util.DeleteFile(imageURL)
		return err
	}
	return nil
}

func (s *ProductService) RemoveImageFromProduct(productID uuid.UUID, imageID uuid.UUID) error {
	image, err := s.productRepository.GetImage(productID, imageID)
	if err != nil {
		return err
	}

	if err := s.productRepository.RemoveImageFromProduct(productID, imageID); err != nil {
		return err
	}
	_ = util.DeleteFile(image.ImageURL)
	return nil
}

func (s *ProductService) AddAttributeToProduct(productID uuid.UUID, attribute *dto.AttributeDto) error {
	if _, err := s.productRepository.GetProduct(productID); err != nil {
		return err
	}

	attributeModel := &models.Attribute{
		ID:        uuid.New(),
		ProductID: productID,
		NameTm:    attribute.NameTm,
		NameRu:    attribute.NameRu,
		NameEn:    attribute.NameEn,
		Type:      attribute.Type,
		Values:    attribute.Values,
	}
	return s.productRepository.AddAttributeToProduct(attributeModel)
}

func (s *ProductService) RemoveAttributeFromProduct(productID uuid.UUID, attributeID uuid.UUID) error {
	return s.productRepository.RemoveAttributeFromProduct(productID, attributeID)
}

// AddValueToAttribute appends a value unless the attribute already has it.
func (s *ProductService) AddValueToAttribute(productID, attributeID uuid.UUID, value string) error {
	attribute, err := s.productRepository.GetAttribute(productID, attributeID)
	if err != nil {
		return err
	}

	if slices.Contains(attribute.Values, value) {
		return nil
	}
	attribute.Values = append(attribute.Values, value)
	return s.productRepository.SaveAttribute(attribute)
}

func (s *ProductService) RemoveValueFromAttribute(productID, attributeID uuid.UUID, value string) error {
	attribute, err := s.productRepository.GetAttribute(productID, attributeID)
	if err != nil {
		return err
	}

	attribute.Values = util.RemoveFromSlice(attribute.Values, value)
	return s.productRepository.SaveAttribute(attribute)
}

func (s *ProductService) GetRecommendedProducts(productID uuid.UUID, paginate *dto.PaginateDto) (*dto.ProductResponseDto, error) {
	paginate.Normalize()

	product, err := s.productRepository.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	products, count, err := s.productRepository.GetRecommendedProducts(paginate, product)
	if err != nil {
		return nil, err
	}

	categoryName := ""
	if product.Category != nil {
		categoryName = product.Category.NameTm
	}

	return &dto.ProductResponseDto{
		Total:        count,
		CategoryName: categoryName,
		Data:         products,
	}, nil
}
