package health

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

func call(t *testing.T, provider DBProvider) (int, Response) {
	t.Helper()
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/health", nil), rec)

	if err := Handler(provider)(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	var body Response
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body %q: %v", rec.Body.String(), err)
	}
	return rec.Code, body
}

func TestHealthUp(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()
	mock.ExpectPing()

	code, body := call(t, func() (*sql.DB, error) { return db, nil })

	if code != http.StatusOK || body.Status != "ok" || body.Database != "up" {
		t.Errorf("got %d %+v, want 200 ok/up", code, body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestHealthDatabaseDown(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()
	mock.ExpectPing().WillReturnError(errors.New("connection refused"))

	code, body := call(t, func() (*sql.DB, error) { return db, nil })

	if code != http.StatusServiceUnavailable || body.Status != "unavailable" || body.Database != "down" {
		t.Errorf("got %d %+v, want 503 unavailable/down", code, body)
	}
}

func TestHealthNoConnectionPool(t *testing.T) {
	code, body := call(t, func() (*sql.DB, error) { return nil, errors.New("no pool") })

	if code != http.StatusServiceUnavailable || body.Database != "down" {
		t.Errorf("got %d %+v, want 503", code, body)
	}
}
