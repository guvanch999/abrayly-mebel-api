package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/validator"
)

func TestBannerDtoValidation(t *testing.T) {
	v := validator.NewValidator()
	valid, invalid := "https://example.com/sale", "not a url"

	if err := v.Validate(&BannerDto{}); err != nil {
		t.Errorf("a banner without link must be valid: %v", err)
	}
	if err := v.Validate(&BannerDto{Link: &valid}); err != nil {
		t.Errorf("valid link rejected: %v", err)
	}
	if err := v.Validate(&BannerDto{Link: &invalid}); err == nil {
		t.Error("invalid link accepted")
	}
}

func TestBannerListToResponse(t *testing.T) {
	link := "https://example.com"
	banners := []*models.Banner{
		{ID: uuid.New(), Image: "uploads/banners/a.jpg", Link: &link},
		{ID: uuid.New(), Image: "uploads/banners/b.jpg"},
	}

	got := BannerListToResponse(banners)

	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	for i, b := range banners {
		if got[i].ID != b.ID || got[i].Image != b.Image || got[i].Link != b.Link {
			t.Errorf("response[%d] = %+v, want fields of %+v", i, got[i], b)
		}
	}
}
