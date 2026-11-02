package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	NameTm    string     `json:"name_tm"`
	NameRu    string     `json:"name_ru"`
	NameEn    string     `json:"name_en"`
	ImageUrl  string     `json:"image_url"`
	ParentId  *uuid.UUID `json:"parent_id"`
	Parent    *Category  `json:"parent" gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Category) TableName() string {
	return "categories"
}
