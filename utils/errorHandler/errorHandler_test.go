package errorHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/labstack/echo/v4"
)

func serve(method string, err error) *httptest.ResponseRecorder {
	e := echo.New()
	req := httptest.NewRequest(method, "/", nil)
	rec := httptest.NewRecorder()
	ResponseHTTPErrorHandler(err, e.NewContext(req, rec))
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) exceptions.ResponseError {
	t.Helper()
	var body exceptions.ResponseError
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not a JSON ResponseError: %v (%q)", err, rec.Body.String())
	}
	return body
}

func TestResponseHTTPErrorHandler(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantCode    int
		wantMessage string
		wantDetails string
	}{
		{
			name:        "ResponseError is returned as is",
			err:         exceptions.NewResponseError(exceptions.ErrNotFound, errors.New("product not found")),
			wantCode:    http.StatusNotFound,
			wantMessage: "Not Found",
			wantDetails: "product not found",
		},
		{
			name:        "wrapped ResponseError is unwrapped",
			err:         fmt.Errorf("service: %w", exceptions.NewResponseError(exceptions.ErrBadRequest, errors.New("bad input"))),
			wantCode:    http.StatusBadRequest,
			wantMessage: "Bad Request",
			wantDetails: "bad input",
		},
		{
			name:        "echo.HTTPError keeps status and message",
			err:         echo.NewHTTPError(http.StatusMethodNotAllowed, "Method Not Allowed"),
			wantCode:    http.StatusMethodNotAllowed,
			wantMessage: "Method Not Allowed",
		},
		{
			name:        "echo route not found",
			err:         echo.ErrNotFound,
			wantCode:    http.StatusNotFound,
			wantMessage: "Not Found",
		},
		{
			name:        "unknown errors are hidden behind a 500",
			err:         errors.New("pq: connection refused to 10.0.0.1"),
			wantCode:    http.StatusInternalServerError,
			wantMessage: "Internal Server Error",
			wantDetails: exceptions.ErrInternalServerError.Details,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(http.MethodGet, tt.err)
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			body := decode(t, rec)
			if body.Code != tt.wantCode {
				t.Errorf("body.code = %d, want %d", body.Code, tt.wantCode)
			}
			if body.Message != tt.wantMessage {
				t.Errorf("body.message = %q, want %q", body.Message, tt.wantMessage)
			}
			if body.Details != tt.wantDetails {
				t.Errorf("body.details = %q, want %q", body.Details, tt.wantDetails)
			}
		})
	}
}

func TestResponseHTTPErrorHandlerHeadHasNoBody(t *testing.T) {
	rec := serve(http.MethodHead, exceptions.ErrUnauthorized)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("HEAD response has a body: %q", rec.Body.String())
	}
}

func TestResponseHTTPErrorHandlerSkipsCommittedResponses(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := c.String(http.StatusOK, "already sent"); err != nil {
		t.Fatalf("write response: %v", err)
	}

	ResponseHTTPErrorHandler(errors.New("late error"), c)

	if rec.Code != http.StatusOK || rec.Body.String() != "already sent" {
		t.Errorf("committed response was modified: %d %q", rec.Code, rec.Body.String())
	}
}
