package models

import "github.com/google/uuid"

type Banner struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;"`
	Image string    `json:"image"`
	Link  *string   `json:"link"`
}

func (b *Banner) TableName() string {
	return "banners"
}
