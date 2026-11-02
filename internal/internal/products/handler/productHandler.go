package handler

import (
	"net/http"

	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/service"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// bindAndValidate binds the request into dst and runs struct validation.
func bindAndValidate(c echo.Context, dst interface{}) error {
	if err := c.Bind(dst); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	if err := c.Validate(dst); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	return nil
}

// ParseListQuery binds and normalizes the pagination and filter query
// parameters of the product list endpoint.
func ParseListQuery(c echo.Context, paginate *dto.PaginateDto, filter *dto.FilterDto) error {
	binder := &echo.DefaultBinder{}
	if err := binder.BindQueryParams(c, paginate); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	if filter != nil {
		if err := binder.BindQueryParams(c, filter); err != nil {
			return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
		}
	}
	paginate.Normalize()
	return nil
}

// @Summary Create product
// @Description Create a product in an existing category
// @Tags admin-products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productDto body dto.ProductDto true "Product"
// @Success 201 {object} util.DataWrapper[models.Product]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError "Category not found"
// @Router /api/v1/admin/products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(dto.ProductDto)
	if err := bindAndValidate(c, product); err != nil {
		return err
	}

	productModel, err := h.productService.CreateProduct(product)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, util.WrapData(productModel))
}

// @Summary Update product
// @Description Update a product by ID
// @Tags admin-products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param productDto body dto.ProductDto true "Product"
// @Success 200 {object} util.DataWrapper[models.Product]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	product := new(dto.ProductDto)
	if err := bindAndValidate(c, product); err != nil {
		return err
	}

	productModel, err := h.productService.UpdateProduct(productID, product)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapData(productModel))
}

// @Summary Delete product
// @Description Delete a product together with its images and attributes
// @Tags admin-products
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := h.productService.DeleteProduct(productID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Product deleted successfully"))
}

// @Summary Get product
// @Description Get a product with its attributes, images and category
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} util.DataWrapper[models.Product]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	product, err := h.productService.GetProduct(productID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapData(product))
}

// @Summary Recommended products
// @Description Other products from the same category, newest first
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 10, max 100)"
// @Success 200 {object} util.DataWrapper[dto.ProductResponseDto]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/products/{id}/recommended [get]
func (h *ProductHandler) GetRecommendedProducts(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	paginate := new(dto.PaginateDto)
	if err := ParseListQuery(c, paginate, nil); err != nil {
		return err
	}

	data, err := h.productService.GetRecommendedProducts(productID, paginate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapData(data))
}

// @Summary List products
// @Description Paginated catalog, newest first. Optional filters: category, discount and a case-insensitive search across the Turkmen, Russian and English names.
// @Tags products
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 10, max 100)"
// @Param category_id query string false "Category ID"
// @Param search query string false "Search term (matches name_tm, name_ru or name_en)"
// @Param is_discount query bool false "true = only discounted products"
// @Success 200 {object} util.DataWrapper[dto.ProductResponseDto]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError "Category not found"
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c echo.Context) error {
	paginate := new(dto.PaginateDto)
	filterDto := new(dto.FilterDto)
	if err := ParseListQuery(c, paginate, filterDto); err != nil {
		return err
	}

	products, err := h.productService.GetProducts(paginate, filterDto)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapData(products))
}

// @Summary Add product image
// @Description Upload an image (.jpg, .jpeg, .png, .webp, .gif, .avif) for a product
// @Tags admin-products
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param image formData file true "Image"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id}/images [post]
func (h *ProductHandler) AddImageToProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	image, err := c.FormFile("image")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := h.productService.AddImageToProduct(productID, image); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Image added to product successfully"))
}

// @Summary Remove product image
// @Description Remove an image from a product and delete the stored file
// @Tags admin-products
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param image_id path string true "Image ID"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id}/images/{image_id} [delete]
func (h *ProductHandler) RemoveImageFromProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	imageID, err := util.GetParamsUUID(c, "image_id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := h.productService.RemoveImageFromProduct(productID, imageID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Image removed from product successfully"))
}

// @Summary Add product attribute
// @Description Add a trilingual attribute (e.g. color, size, material) with its values
// @Tags admin-products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param attributeDto body dto.AttributeDto true "Attribute"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id}/attributes [post]
func (h *ProductHandler) AddAttributeToProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	attribute := new(dto.AttributeDto)
	if err := bindAndValidate(c, attribute); err != nil {
		return err
	}

	if err := h.productService.AddAttributeToProduct(productID, attribute); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Attribute added to product successfully"))
}

// @Summary Remove product attribute
// @Description Remove an attribute from a product
// @Tags admin-products
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param attribute_id path string true "Attribute ID"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id}/attributes/{attribute_id} [delete]
func (h *ProductHandler) RemoveAttributeFromProduct(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	attributeID, err := util.GetParamsUUID(c, "attribute_id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := h.productService.RemoveAttributeFromProduct(productID, attributeID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Attribute removed from product successfully"))
}

// @Summary Add attribute value
// @Description Append a value to a product attribute (no-op if it already exists)
// @Tags admin-products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param attribute_id path string true "Attribute ID"
// @Param body body dto.AttributeValueDto true "Value"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id}/attributes/{attribute_id}/values [post]
func (h *ProductHandler) AddValueToAttribute(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	valueDto := new(dto.AttributeValueDto)
	if err := bindAndValidate(c, valueDto); err != nil {
		return err
	}

	if err := h.productService.AddValueToAttribute(productID, valueDto.AttributeID, valueDto.Value); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Value added to attribute successfully"))
}

// @Summary Remove attribute value
// @Description Remove a value from a product attribute
// @Tags admin-products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param attribute_id path string true "Attribute ID"
// @Param body body dto.AttributeValueDto true "Value"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/products/{id}/attributes/{attribute_id}/values/delete [post]
func (h *ProductHandler) RemoveValueFromAttribute(c echo.Context) error {
	productID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	valueDto := new(dto.AttributeValueDto)
	if err := bindAndValidate(c, valueDto); err != nil {
		return err
	}

	if err := h.productService.RemoveValueFromAttribute(productID, valueDto.AttributeID, valueDto.Value); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Value removed from attribute successfully"))
}
