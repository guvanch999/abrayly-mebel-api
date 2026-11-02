package dto

import (
	"strings"

	"github.com/google/uuid"
)

// FilterDto holds the optional catalog filters read from the query string.
type FilterDto struct {
	CategoryID *uuid.UUID `json:"category_id" form:"category_id" query:"category_id"`
	Search     *string    `json:"search" form:"search" query:"search"`
	IsDiscount *bool      `json:"is_discount" form:"is_discount" query:"is_discount"`
}

// SearchTerm returns the trimmed search string, or "" when none was given.
func (f *FilterDto) SearchTerm() string {
	if f == nil || f.Search == nil {
		return ""
	}
	return strings.TrimSpace(*f.Search)
}

// OnlyDiscounted reports whether only discounted products were requested
// (is_discount=true). is_discount=false or no value applies no filter.
func (f *FilterDto) OnlyDiscounted() bool {
	return f != nil && f.IsDiscount != nil && *f.IsDiscount
}
