package handler

import (
	"net/http"

	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/auth/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/auth/service"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	adminService *service.AdminService
}

func NewAuthHandler(adminService *service.AdminService) *AuthHandler {
	return &AuthHandler{adminService: adminService}
}

// @Summary Admin login
// @Description Exchange admin credentials for an access token (30 min) and a refresh token (24 h)
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequestDto body dto.LoginRequestDto true "Login request"
// @Success 200 {object} util.DataWrapper[dto.TokenResponseDto]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Router /auth/admin/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var loginRequestDto dto.LoginRequestDto
	if err := c.Bind(&loginRequestDto); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := c.Validate(&loginRequestDto); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	token, err := h.adminService.Login(loginRequestDto.Username, loginRequestDto.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, util.WrapData(token))
}

// @Summary Refresh tokens
// @Description Exchange a valid refresh token for a new token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshTokenRequestDto body dto.RefreshTokenRequestDto true "Refresh token request"
// @Success 200 {object} util.DataWrapper[dto.TokenResponseDto]
// @Failure 400 {object} exceptions.ResponseError
// @Failure 401 {object} exceptions.ResponseError
// @Router /auth/admin/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var refreshTokenRequestDto dto.RefreshTokenRequestDto
	if err := c.Bind(&refreshTokenRequestDto); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	if err := c.Validate(&refreshTokenRequestDto); err != nil {
		return exceptions.NewResponseError(exceptions.ErrBadRequest, err)
	}

	token, err := h.adminService.RefreshToken(refreshTokenRequestDto.RefreshToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, util.WrapData(token))
}
