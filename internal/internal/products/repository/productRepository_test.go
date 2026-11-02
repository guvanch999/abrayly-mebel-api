package repository

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return db, mock
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// dryRun renders the SELECT produced by FilterScope without a database.
func dryRun(t *testing.T, filter *dto.FilterDto) (string, []interface{}) {
	t.Helper()
	db, _ := newMockDB(t)
	var products []models.Product
	stmt := db.Session(&gorm.Session{DryRun: true}).
		Model(&models.Product{}).
		Scopes(FilterScope(filter)).
		Find(&products).Statement
	return stmt.SQL.String(), stmt.Vars
}

func TestFilterScopeBuildsQuery(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name      string
		filter    *dto.FilterDto
		wantSQL   []string
		forbidSQL []string
		wantVars  []interface{}
	}{
		{
			name:      "no filters",
			filter:    &dto.FilterDto{},
			wantSQL:   []string{`SELECT * FROM "products"`},
			forbidSQL: []string{"WHERE"},
		},
		{
			name:     "category",
			filter:   &dto.FilterDto{CategoryID: &categoryID},
			wantSQL:  []string{"WHERE category_id = $1"},
			wantVars: []interface{}{categoryID},
		},
		{
			name:     "trilingual search",
			filter:   &dto.FilterDto{Search: strPtr("  Sofa ")},
			wantSQL:  []string{"name_tm ILIKE $1 OR name_ru ILIKE $2 OR name_en ILIKE $3"},
			wantVars: []interface{}{"%Sofa%", "%Sofa%", "%Sofa%"},
		},
		{
			name:      "blank search is ignored",
			filter:    &dto.FilterDto{Search: strPtr("   ")},
			forbidSQL: []string{"ILIKE"},
		},
		{
			name:     "search escapes LIKE wildcards",
			filter:   &dto.FilterDto{Search: strPtr("50%_off")},
			wantVars: []interface{}{`%50\%\_off%`, `%50\%\_off%`, `%50\%\_off%`},
		},
		{
			name:    "only discounted",
			filter:  &dto.FilterDto{IsDiscount: boolPtr(true)},
			wantSQL: []string{"WHERE discount > 0"},
		},
		{
			name:      "is_discount=false does not filter",
			filter:    &dto.FilterDto{IsDiscount: boolPtr(false)},
			forbidSQL: []string{"discount"},
		},
		{
			name:     "all filters are ANDed and the search is parenthesized",
			filter:   &dto.FilterDto{CategoryID: &categoryID, Search: strPtr("diwan"), IsDiscount: boolPtr(true)},
			wantSQL:  []string{"WHERE category_id = $1 AND (name_tm ILIKE $2 OR name_ru ILIKE $3 OR name_en ILIKE $4) AND discount > 0"},
			wantVars: []interface{}{categoryID, "%diwan%", "%diwan%", "%diwan%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, vars := dryRun(t, tt.filter)
			for _, want := range tt.wantSQL {
				if !strings.Contains(sql, want) {
					t.Errorf("SQL %q does not contain %q", sql, want)
				}
			}
			for _, forbidden := range tt.forbidSQL {
				if strings.Contains(sql, forbidden) {
					t.Errorf("SQL %q must not contain %q", sql, forbidden)
				}
			}
			if tt.wantVars != nil && !reflect.DeepEqual(vars, tt.wantVars) {
				t.Errorf("vars = %#v, want %#v", vars, tt.wantVars)
			}
		})
	}
}

func TestEscapeLike(t *testing.T) {
	tests := map[string]string{
		"sofa":       "sofa",
		"100%":       `100\%`,
		"a_b":        `a\_b`,
		`back\slash`: `back\\slash`,
	}
	for in, want := range tests {
		if got := EscapeLike(in); got != want {
			t.Errorf("EscapeLike(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestGetProductsCountsThenLoadsPageWithAssociations(t *testing.T) {
	db, mock := newMockDB(t)
	mock.MatchExpectationsInOrder(false)
	repo := NewProductRepository(db)

	categoryID := uuid.New()
	productID := uuid.New()
	now := time.Now().UTC()
	filter := &dto.FilterDto{CategoryID: &categoryID, Search: strPtr("sofa"), IsDiscount: boolPtr(true)}
	paginate := &dto.PaginateDto{Page: 3, PageSize: 5}

	where := regexp.QuoteMeta(`WHERE category_id = $1 AND (name_tm ILIKE $2 OR name_ru ILIKE $3 OR name_en ILIKE $4) AND discount > 0`)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" `)+where).
		WithArgs(categoryID, "%sofa%", "%sofa%", "%sofa%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(11))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" `)+where+regexp.QuoteMeta(` ORDER BY created_at DESC, id LIMIT $5 OFFSET $6`)).
		WithArgs(categoryID, "%sofa%", "%sofa%", "%sofa%", 5, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "name_tm", "name_ru", "name_en", "price", "discount", "old_price", "created_at", "updated_at"}).
			AddRow(productID.String(), categoryID.String(), "Diwan", "Диван", "Sofa", 1000.0, 10.0, 1100.0, now, now))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM "attributes"`)).
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name_tm", "name_ru", "name_en", "type", "values", "product_id"}).
			AddRow(uuid.NewString(), "Reňk", "Цвет", "Color", "color", `["Red","Blue"]`, productID.String()))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM "product_images"`)).
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "image_url"}).
			AddRow(uuid.NewString(), productID.String(), "uploads/products/1.jpg"))

	products, total, err := repo.GetProducts(paginate, filter)
	if err != nil {
		t.Fatalf("GetProducts: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	if total != 11 {
		t.Errorf("total = %d, want 11", total)
	}
	if len(products) != 1 || products[0].ID != productID {
		t.Fatalf("products = %+v, want one product %s", products, productID)
	}
	if got := products[0].Attributes; len(got) != 1 || !reflect.DeepEqual(got[0].Values, []string{"Red", "Blue"}) {
		t.Errorf("attributes = %+v, want JSON values [Red Blue]", got)
	}
	if got := products[0].Images; len(got) != 1 || got[0].ImageURL != "uploads/products/1.jpg" {
		t.Errorf("images = %+v", got)
	}
}

func TestGetProductsSkipsPageQueryWhenNothingMatches(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	products, total, err := repo.GetProducts(&dto.PaginateDto{Page: 1, PageSize: 10}, &dto.FilterDto{})
	if err != nil {
		t.Fatalf("GetProducts: %v", err)
	}
	if total != 0 || products == nil || len(products) != 0 {
		t.Errorf("got total=%d products=%#v, want 0 and an empty (non-nil) slice", total, products)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProductRemovesChildrenInOneTransaction(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)
	productID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "product_images" WHERE product_id = $1`)).
		WithArgs(productID).WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "attributes" WHERE product_id = $1`)).
		WithArgs(productID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "products" WHERE id = $1`)).
		WithArgs(productID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.DeleteProduct(productID); err != nil {
		t.Fatalf("DeleteProduct: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProductRollsBackOnError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)
	productID := uuid.New()
	boom := errors.New("boom")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "product_images"`)).
		WithArgs(productID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "attributes"`)).
		WithArgs(productID).WillReturnError(boom)
	mock.ExpectRollback()

	if err := repo.DeleteProduct(productID); !errors.Is(err, boom) {
		t.Fatalf("DeleteProduct error = %v, want %v", err, boom)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
