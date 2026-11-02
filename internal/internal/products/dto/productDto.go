package dto

import "github.com/google/uuid"

// ProductDto is the body for creating or updating a product.
type ProductDto struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required,uuid"`
	NameTm     string    `json:"name_tm" validate:"required"`
	NameRu     string    `json:"name_ru" validate:"required"`
	NameEn     string    `json:"name_en" validate:"required"`
	Price      float64   `json:"price" validate:"required,gt=0"`
	Discount   float64   `json:"discount" validate:"gte=0"`
	OldPrice   float64   `json:"old_price" validate:"gte=0"`
}

// AttributeDto is the body for adding an attribute (e.g. color, size) to a product.
type AttributeDto struct {
	NameTm string   `json:"name_tm" validate:"required"`
	NameRu string   `json:"name_ru" validate:"required"`
	NameEn string   `json:"name_en" validate:"required"`
	Type   string   `json:"type" validate:"required"`
	Values []string `json:"values" validate:"required,min=1,dive,required"`
}

// AttributeValueDto adds or removes one value of an attribute.
type AttributeValueDto struct {
	AttributeID uuid.UUID `query:"attribute_id" param:"attribute_id" validate:"required,uuid"`
	Value       string    `json:"value" validate:"required"`
}
