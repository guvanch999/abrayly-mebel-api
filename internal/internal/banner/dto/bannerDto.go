package dto

import (
	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
)

type BannerDto struct {
	Link *string `json:"link" form:"link" validate:"omitempty,url"`
}

type BannerResponse struct {
	ID    uuid.UUID `json:"id"`
	Link  *string   `json:"link"`
	Image string    `json:"image"`
}

func BannerToResponse(banner *models.Banner) *BannerResponse {
	return &BannerResponse{
		ID:    banner.ID,
		Link:  banner.Link,
		Image: banner.Image,
	}
}

func BannerListToResponse(banners []*models.Banner) []*BannerResponse {
	response := make([]*BannerResponse, len(banners))
	for i, banner := range banners {
		response[i] = BannerToResponse(banner)
	}
	return response
}
