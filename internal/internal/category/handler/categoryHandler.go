package handler

import (
	"net/http"

	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/service"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func bindCategory(c echo.Context) (*dto.CategoryDto, error) {
	categoryDto := new(dto.CategoryDto)
	if err := c.Bind(categoryDto); err != nil {
		return nil, exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	if err := c.Validate(categoryDto); err != nil {
		return nil, exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	return categoryDto, nil
}

// @Summary Create category
// @Description Create a root category or, with parent_id, a subcategory
// @Tags admin-categories
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param name_tm formData string true "Name in Turkmen"
// @Param name_ru formData string true "Name in Russian"
// @Param name_en formData string true "Name in English"
// @Param parent_id formData string false "Parent category ID"
// @Param image formData file true "Image (.jpg, .jpeg, .png, .webp, .gif, .avif)"
// @Success 201 {object} util.DataWrapper[models.Category]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError "Parent category not found"
// @Router /api/v1/admin/categories [post]
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	categoryDto, err := bindCategory(c)
	if err != nil {
		return err
	}

	file, err := c.FormFile("image")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	category, err := h.service.CreateCategory(categoryDto, file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, util.WrapData(category))
}

// @Summary Get category
// @Description Get a category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} util.DataWrapper[models.Category]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	category, err := h.service.GetCategory(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapData(category))
}

// @Summary Category tree
// @Description All categories as a tree: root categories with nested subcategories
// @Tags categories
// @Produce json
// @Success 200 {object} util.ListDataWrapper[[]dto.CategoryResponseDto]
// @Failure 500 {object} exceptions.ResponseError
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapListData(categories))
}

// @Summary Update category
// @Description Update names, parent and (optionally) the image of a category
// @Tags admin-categories
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID"
// @Param name_tm formData string true "Name in Turkmen"
// @Param name_ru formData string true "Name in Russian"
// @Param name_en formData string true "Name in English"
// @Param parent_id formData string false "Parent category ID"
// @Param image formData file false "New image"
// @Success 200 {object} util.DataWrapper[models.Category]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	categoryDto, err := bindCategory(c)
	if err != nil {
		return err
	}

	// The image is optional on update.
	file, _ := c.FormFile("image")

	category, err := h.service.UpdateCategory(id, categoryDto, file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapData(category))
}

// @Summary Delete category
// @Description Delete a category; subcategories are deleted by cascade, products are kept with category_id = null
// @Tags admin-categories
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /api/v1/admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := h.service.DeleteCategory(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, util.WrapMessage("Category deleted successfully"))
}
