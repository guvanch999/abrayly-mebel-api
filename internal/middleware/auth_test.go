package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/errorHandler"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"github.com/labstack/echo/v4"
)

func TestBearerToken(t *testing.T) {
	tests := []struct {
		header  string
		want    string
		wantErr bool
	}{
		{header: "Bearer abc.def.ghi", want: "abc.def.ghi"},
		{header: "bearer abc.def.ghi", want: "abc.def.ghi"},
		{header: "  Bearer   abc.def.ghi  ", want: "abc.def.ghi"},
		{header: "", wantErr: true},
		{header: "abc.def.ghi", wantErr: true},
		{header: "Basic dXNlcjpwYXNz", wantErr: true},
		{header: "Bearer ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			got, err := BearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BearerToken(%q) error = %v, wantErr %v", tt.header, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("BearerToken(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

type fakeFinder map[uuid.UUID]*models.Admin

func (f fakeFinder) FindByID(id uuid.UUID) (*models.Admin, error) {
	if admin, ok := f[id]; ok {
		return admin, nil
	}
	return nil, errors.New("record not found")
}

const (
	accessSecret  = "access-secret"
	refreshSecret = "refresh-secret"
)

func newProtectedServer(admins fakeFinder) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = errorHandler.ResponseHTTPErrorHandler
	e.GET("/admin/ping", func(c echo.Context) error {
		admin := c.Get(AdminContextKey).(*models.Admin)
		return c.String(http.StatusOK, admin.Username)
	}, RequireAdmin(accessSecret, admins))
	return e
}

func issue(t *testing.T, id uuid.UUID, secret string, ttl time.Duration) string {
	t.Helper()
	token, err := util.IssueToken(id, secret, ttl)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	return token
}

func TestRequireAdmin(t *testing.T) {
	admin := &models.Admin{ID: uuid.New(), Username: "owner"}
	e := newProtectedServer(fakeFinder{admin.ID: admin})

	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantBody   string
	}{
		{name: "valid access token", header: "Bearer " + issue(t, admin.ID, accessSecret, time.Minute), wantStatus: http.StatusOK, wantBody: "owner"},
		{name: "missing header", header: "", wantStatus: http.StatusUnauthorized},
		{name: "malformed header", header: "Token xyz", wantStatus: http.StatusUnauthorized},
		{name: "refresh token used as access token", header: "Bearer " + issue(t, admin.ID, refreshSecret, time.Minute), wantStatus: http.StatusUnauthorized},
		{name: "expired token", header: "Bearer " + issue(t, admin.ID, accessSecret, -time.Minute), wantStatus: http.StatusUnauthorized},
		{name: "unknown admin", header: "Bearer " + issue(t, uuid.New(), accessSecret, time.Minute), wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/ping", nil)
			if tt.header != "" {
				req.Header.Set(echo.HeaderAuthorization, tt.header)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
			if tt.wantStatus == http.StatusUnauthorized {
				var body map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body["message"] != "Unauthorized" {
					t.Errorf("unexpected 401 body: %s", rec.Body.String())
				}
			}
		})
	}
}
