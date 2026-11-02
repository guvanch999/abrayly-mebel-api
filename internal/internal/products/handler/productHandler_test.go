package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/dto"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/labstack/echo/v4"
)

func parse(t *testing.T, query string) (*dto.PaginateDto, *dto.FilterDto, error) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?"+query, nil)
	c := e.NewContext(req, httptest.NewRecorder())

	paginate, filter := new(dto.PaginateDto), new(dto.FilterDto)
	err := ParseListQuery(c, paginate, filter)
	return paginate, filter, err
}

func TestParseListQueryDefaults(t *testing.T) {
	paginate, filter, err := parse(t, "")
	if err != nil {
		t.Fatalf("ParseListQuery: %v", err)
	}
	if paginate.Page != dto.DefaultPage || paginate.PageSize != dto.DefaultPageSize {
		t.Errorf("pagination = %+v, want defaults", *paginate)
	}
	if filter.CategoryID != nil || filter.Search != nil || filter.IsDiscount != nil {
		t.Errorf("filter = %+v, want no filters", *filter)
	}
}

func TestParseListQueryReadsAllParameters(t *testing.T) {
	categoryID := uuid.New()

	paginate, filter, err := parse(t, "page=3&limit=25&category_id="+categoryID.String()+"&search=%D0%B4%D0%B8%D0%B2%D0%B0%D0%BD&is_discount=true")
	if err != nil {
		t.Fatalf("ParseListQuery: %v", err)
	}

	if paginate.Page != 3 || paginate.PageSize != 25 {
		t.Errorf("pagination = %+v, want page 3, limit 25", *paginate)
	}
	if filter.CategoryID == nil || *filter.CategoryID != categoryID {
		t.Errorf("category_id = %v, want %s", filter.CategoryID, categoryID)
	}
	if filter.SearchTerm() != "диван" {
		t.Errorf("search = %q, want %q (Cyrillic, URL-encoded)", filter.SearchTerm(), "диван")
	}
	if !filter.OnlyDiscounted() {
		t.Error("is_discount=true was not parsed")
	}
}

func TestParseListQueryCapsLimit(t *testing.T) {
	paginate, _, err := parse(t, "page=0&limit=100000")
	if err != nil {
		t.Fatalf("ParseListQuery: %v", err)
	}
	if paginate.Page != 1 || paginate.PageSize != dto.MaxPageSize {
		t.Errorf("pagination = %+v, want page 1, limit %d", *paginate, dto.MaxPageSize)
	}
}

func TestParseListQueryRejectsInvalidValues(t *testing.T) {
	for _, query := range []string{
		"page=abc",
		"limit=ten",
		"category_id=not-a-uuid",
		"is_discount=maybe",
	} {
		t.Run(query, func(t *testing.T) {
			_, _, err := parse(t, query)
			var respErr *exceptions.ResponseError
			if !errors.As(err, &respErr) || respErr.Code != http.StatusBadRequest {
				t.Fatalf("error = %v, want a 400 ResponseError", err)
			}
		})
	}
}
