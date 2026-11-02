package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;default:null"`
	NameTm     string    `json:"name_tm"`
	NameRu     string    `json:"name_ru"`
	NameEn     string    `json:"name_en"`

	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
	OldPrice float64 `json:"old_price"`

	Attributes []Attribute    `json:"attributes" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Images     []ProductImage `json:"images" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Category   *Category      `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Product) TableName() string {
	return "products"
}
