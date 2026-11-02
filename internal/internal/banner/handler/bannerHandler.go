package handler

import (
	"net/http"

	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/service"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"github.com/labstack/echo/v4"
)

type BannerHandler struct {
	bannerService *service.BannerService
}

func NewBannerHandler(bannerService *service.BannerService) *BannerHandler {
	return &BannerHandler{bannerService: bannerService}
}

func bindBanner(c echo.Context) (*dto.BannerDto, error) {
	bannerDto := new(dto.BannerDto)
	if err := c.Bind(bannerDto); err != nil {
		return nil, exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	if err := c.Validate(bannerDto); err != nil {
		return nil, exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	return bannerDto, nil
}

// @Summary Create banner
// @Description Upload a home-page banner image with an optional link
// @Tags admin-banners
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param link formData string false "Link opened when the banner is clicked"
// @Param image formData file true "Image (.jpg, .jpeg, .png, .webp, .gif, .avif)"
// @Success 201 {object} util.DataWrapper[dto.BannerResponse]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Router /admin/banners [post]
func (h *BannerHandler) Create(c echo.Context) error {
	bannerDto, err := bindBanner(c)
	if err != nil {
		return err
	}
	image, err := c.FormFile("image")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	createdBanner, err := h.bannerService.Create(bannerDto.Link, image)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, util.WrapData(createdBanner))
}

// @Summary List banners
// @Description Get all banners
// @Tags banners
// @Produce json
// @Success 200 {object} util.ListDataWrapper[[]dto.BannerResponse]
// @Failure 500 {object} exceptions.ResponseError
// @Router /banners [get]
func (h *BannerHandler) GetAll(c echo.Context) error {
	banners, err := h.bannerService.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, util.WrapListData(banners))
}

// @Summary Get banner
// @Description Get a banner by ID
// @Tags banners
// @Produce json
// @Param id path string true "Banner ID"
// @Success 200 {object} util.DataWrapper[dto.BannerResponse]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /banners/{id} [get]
func (h *BannerHandler) GetByID(c echo.Context) error {
	bannerUUID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	banner, err := h.bannerService.GetByID(bannerUUID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, util.WrapData(banner))
}

// @Summary Update banner
// @Description Change the link of a banner
// @Tags admin-banners
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Banner ID"
// @Param bannerDto body dto.BannerDto true "Banner"
// @Success 200 {object} util.DataWrapper[dto.BannerResponse]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /admin/banners/{id} [put]
func (h *BannerHandler) Update(c echo.Context) error {
	bannerUUID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	bannerDto, err := bindBanner(c)
	if err != nil {
		return err
	}
	updatedBanner, err := h.bannerService.Update(bannerUUID, bannerDto.Link)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, util.WrapData(updatedBanner))
}

// @Summary Delete banner
// @Description Delete a banner and its stored image
// @Tags admin-banners
// @Produce json
// @Security Bearer
// @Param id path string true "Banner ID"
// @Success 200 {object} util.MessageWrapper
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Failure 404 {object} exceptions.ResponseError
// @Router /admin/banners/{id} [delete]
func (h *BannerHandler) Delete(c echo.Context) error {
	bannerUUID, err := util.GetParamsUUID(c, "id")
	if err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}
	if err := h.bannerService.Delete(bannerUUID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, util.WrapMessage("Banner deleted"))
}
