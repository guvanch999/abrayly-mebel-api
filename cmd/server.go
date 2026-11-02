package cmd

import (
	"github.com/guvanch999/abrayly-mebel-backend/internal/config"
	"github.com/guvanch999/abrayly-mebel-backend/internal/database"
	"github.com/guvanch999/abrayly-mebel-backend/utils/errorHandler"
	"github.com/guvanch999/abrayly-mebel-backend/utils/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	Echo     *echo.Echo
	Config   *config.Schema
	Database *gorm.DB
}

func NewServer() *Server {
	cfg := config.GetConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = errorHandler.ResponseHTTPErrorHandler

	return &Server{
		Echo:     e,
		Config:   cfg,
		Database: db,
	}
}
