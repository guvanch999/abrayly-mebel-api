package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/utils/validator"
)

func TestPaginateNormalize(t *testing.T) {
	tests := []struct {
		name       string
		in         PaginateDto
		want       PaginateDto
		wantOffset int
	}{
		{name: "defaults", in: PaginateDto{}, want: PaginateDto{Page: 1, PageSize: 10}, wantOffset: 0},
		{name: "keeps valid values", in: PaginateDto{Page: 3, PageSize: 20}, want: PaginateDto{Page: 3, PageSize: 20}, wantOffset: 40},
		{name: "negative values", in: PaginateDto{Page: -2, PageSize: -5}, want: PaginateDto{Page: 1, PageSize: 10}, wantOffset: 0},
		{name: "caps page size", in: PaginateDto{Page: 2, PageSize: 10_000}, want: PaginateDto{Page: 2, PageSize: MaxPageSize}, wantOffset: MaxPageSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in
			got.Normalize()
			if got != tt.want {
				t.Errorf("Normalize() = %+v, want %+v", got, tt.want)
			}
			if offset := got.Offset(); offset != tt.wantOffset {
				t.Errorf("Offset() = %d, want %d", offset, tt.wantOffset)
			}
		})
	}
}

func TestFilterHelpers(t *testing.T) {
	yes, no := true, false
	blank, term := "   ", "  sofa  "

	if (&FilterDto{}).OnlyDiscounted() {
		t.Error("no is_discount must not filter")
	}
	if (&FilterDto{IsDiscount: &no}).OnlyDiscounted() {
		t.Error("is_discount=false must not filter")
	}
	if !(&FilterDto{IsDiscount: &yes}).OnlyDiscounted() {
		t.Error("is_discount=true must filter")
	}

	if got := (&FilterDto{}).SearchTerm(); got != "" {
		t.Errorf("SearchTerm() = %q, want empty", got)
	}
	if got := (&FilterDto{Search: &blank}).SearchTerm(); got != "" {
		t.Errorf("SearchTerm() = %q, want empty for blank input", got)
	}
	if got := (&FilterDto{Search: &term}).SearchTerm(); got != "sofa" {
		t.Errorf("SearchTerm() = %q, want %q", got, "sofa")
	}

	var nilFilter *FilterDto
	if nilFilter.OnlyDiscounted() || nilFilter.SearchTerm() != "" {
		t.Error("a nil filter must behave like an empty one")
	}
}

func validProduct() ProductDto {
	return ProductDto{
		CategoryID: uuid.New(),
		NameTm:     "Diwan",
		NameRu:     "Диван",
		NameEn:     "Sofa",
		Price:      1000,
		Discount:   10,
		OldPrice:   1100,
	}
}

func TestProductDtoValidation(t *testing.T) {
	v := validator.NewValidator()

	tests := []struct {
		name    string
		mutate  func(p *ProductDto)
		wantErr bool
	}{
		{name: "valid", mutate: func(p *ProductDto) {}},
		{name: "no discount is fine", mutate: func(p *ProductDto) { p.Discount, p.OldPrice = 0, 0 }},
		{name: "missing category", mutate: func(p *ProductDto) { p.CategoryID = uuid.Nil }, wantErr: true},
		{name: "missing Turkmen name", mutate: func(p *ProductDto) { p.NameTm = "" }, wantErr: true},
		{name: "missing Russian name", mutate: func(p *ProductDto) { p.NameRu = "" }, wantErr: true},
		{name: "missing English name", mutate: func(p *ProductDto) { p.NameEn = "" }, wantErr: true},
		{name: "zero price", mutate: func(p *ProductDto) { p.Price = 0 }, wantErr: true},
		{name: "negative price", mutate: func(p *ProductDto) { p.Price = -1 }, wantErr: true},
		{name: "negative discount", mutate: func(p *ProductDto) { p.Discount = -5 }, wantErr: true},
		{name: "negative old price", mutate: func(p *ProductDto) { p.OldPrice = -5 }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validProduct()
			tt.mutate(&p)
			err := v.Validate(&p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttributeDtoValidation(t *testing.T) {
	v := validator.NewValidator()
	valid := func() AttributeDto {
		return AttributeDto{NameTm: "Reňk", NameRu: "Цвет", NameEn: "Color", Type: "color", Values: []string{"Red"}}
	}

	tests := []struct {
		name    string
		mutate  func(a *AttributeDto)
		wantErr bool
	}{
		{name: "valid", mutate: func(a *AttributeDto) {}},
		{name: "missing type", mutate: func(a *AttributeDto) { a.Type = "" }, wantErr: true},
		{name: "missing name", mutate: func(a *AttributeDto) { a.NameEn = "" }, wantErr: true},
		{name: "no values", mutate: func(a *AttributeDto) { a.Values = nil }, wantErr: true},
		{name: "empty values", mutate: func(a *AttributeDto) { a.Values = []string{} }, wantErr: true},
		{name: "blank value", mutate: func(a *AttributeDto) { a.Values = []string{"Red", ""} }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := valid()
			tt.mutate(&a)
			err := v.Validate(&a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttributeValueDtoValidation(t *testing.T) {
	v := validator.NewValidator()

	if err := v.Validate(&AttributeValueDto{AttributeID: uuid.New(), Value: "Blue"}); err != nil {
		t.Errorf("valid value rejected: %v", err)
	}
	if err := v.Validate(&AttributeValueDto{AttributeID: uuid.Nil, Value: "Blue"}); err == nil {
		t.Error("missing attribute ID accepted")
	}
	if err := v.Validate(&AttributeValueDto{AttributeID: uuid.New()}); err == nil {
		t.Error("empty value accepted")
	}
}
