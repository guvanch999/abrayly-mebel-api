package models

import "github.com/google/uuid"

type Attribute struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	NameTm    string    `json:"name_tm" gorm:"type:varchar(255);not null"`
	NameRu    string    `json:"name_ru" gorm:"type:varchar(255);not null"`
	NameEn    string    `json:"name_en" gorm:"type:varchar(255);not null"`
	Type      string    `json:"type" gorm:"type:varchar(50);not null"`
	Values    []string  `json:"values" gorm:"type:text;serializer:json"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid"`
}

func (a *Attribute) TableName() string {
	return "attributes"
}
