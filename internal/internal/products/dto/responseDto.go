package dto

import (
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
)

type ProductResponseDto struct {
	Total        int              `json:"total"`
	CategoryName string           `json:"category_name"`
	Data         []models.Product `json:"data"`
}
