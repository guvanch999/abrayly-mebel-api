package dto

const (
	// DefaultPage is used when the page query parameter is missing or invalid.
	DefaultPage = 1
	// DefaultPageSize is used when the limit query parameter is missing or invalid.
	DefaultPageSize = 10
	// MaxPageSize caps the limit query parameter.
	MaxPageSize = 100
)

// PaginateDto holds 1-based page/limit query parameters.
type PaginateDto struct {
	Page     int `query:"page"`
	PageSize int `query:"limit"`
}

// Normalize replaces missing or out-of-range values with defaults and caps
// the page size at MaxPageSize.
func (p *PaginateDto) Normalize() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}

// Offset returns the number of rows to skip for the current page.
func (p *PaginateDto) Offset() int {
	return (p.Page - 1) * p.PageSize
}
