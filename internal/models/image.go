package models

import "github.com/google/uuid"

type ProductImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid"`
	ImageURL  string    `json:"image_url"`
}

func (p *ProductImage) TableName() string {
	return "product_images"
}
