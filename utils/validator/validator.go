package validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// StructValidator adapts go-playground/validator to echo.Validator.
type StructValidator struct {
	validator *validator.Validate
}

// Validate validates struct tags and maps failures to HTTP 400.
func (cv *StructValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// NewValidator returns a validator ready to be assigned to echo.Echo.Validator.
func NewValidator() *StructValidator {
	return &StructValidator{
		validator: validator.New(),
	}
}
