package models

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"type:varchar(1000);not null"`
	Username  string    `gorm:"type:varchar(1000);not null;unique"`
	Password  string    `gorm:"type:varchar(1000);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Admin) TableName() string {
	return "admins"
}
