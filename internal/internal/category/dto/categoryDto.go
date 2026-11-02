package dto

import "github.com/google/uuid"

type CategoryDto struct {
	NameTm   string     `json:"name_tm" form:"name_tm" validate:"required"`
	NameRu   string     `json:"name_ru" form:"name_ru" validate:"required"`
	NameEn   string     `json:"name_en" form:"name_en" validate:"required"`
	ParentId *uuid.UUID `json:"parent_id" form:"parent_id" validate:"omitempty,uuid"`
}
