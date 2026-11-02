package util

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
)

const testSecret = "test-access-secret"

func TestIssueAndCheckTokenRoundTrip(t *testing.T) {
	id := uuid.New()

	token, err := IssueToken(id, testSecret, time.Minute)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	got, err := CheckToken(token, testSecret)
	if err != nil {
		t.Fatalf("CheckToken: %v", err)
	}
	if got != id {
		t.Fatalf("subject = %s, want %s", got, id)
	}
}

func TestIssueTokenRequiresSecret(t *testing.T) {
	if _, err := IssueToken(uuid.New(), "", time.Minute); err == nil {
		t.Fatal("expected an error for an empty secret")
	}
}

func TestCheckTokenRejectsInvalidTokens(t *testing.T) {
	id := uuid.New()
	future := jwt.NewNumericDate(time.Now().Add(time.Hour))

	sign := func(method jwt.SigningMethod, key interface{}, claims jwt.Claims) string {
		t.Helper()
		token, err := jwt.NewWithClaims(method, claims).SignedString(key)
		if err != nil {
			t.Fatalf("sign: %v", err)
		}
		return token
	}

	valid, err := IssueToken(id, testSecret, time.Minute)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	expired, err := IssueToken(id, testSecret, -time.Minute)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		secret  string
		details string
	}{
		{name: "empty token", token: "", secret: testSecret, details: "token is required"},
		{name: "malformed token", token: "not.a.jwt", secret: testSecret, details: "invalid token"},
		{name: "wrong secret", token: valid, secret: "another-secret", details: "invalid token"},
		{name: "expired", token: expired, secret: testSecret, details: "token expired"},
		{
			name:    "unexpected algorithm HS512",
			token:   sign(jwt.SigningMethodHS512, []byte(testSecret), jwt.RegisteredClaims{Subject: id.String(), ExpiresAt: future}),
			secret:  testSecret,
			details: "invalid token",
		},
		{
			name:    "alg none",
			token:   sign(jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, jwt.RegisteredClaims{Subject: id.String(), ExpiresAt: future}),
			secret:  testSecret,
			details: "invalid token",
		},
		{
			name:    "missing expiry",
			token:   sign(jwt.SigningMethodHS256, []byte(testSecret), jwt.RegisteredClaims{Subject: id.String()}),
			secret:  testSecret,
			details: "invalid token",
		},
		{
			name:    "subject is not a UUID",
			token:   sign(jwt.SigningMethodHS256, []byte(testSecret), jwt.RegisteredClaims{Subject: "42", ExpiresAt: future}),
			secret:  testSecret,
			details: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckToken(tt.token, tt.secret)
			if got != uuid.Nil {
				t.Errorf("subject = %s, want uuid.Nil", got)
			}

			var respErr *exceptions.ResponseError
			if !errors.As(err, &respErr) {
				t.Fatalf("error = %v (%T), want *exceptions.ResponseError", err, err)
			}
			if respErr.Code != http.StatusUnauthorized {
				t.Errorf("code = %d, want %d", respErr.Code, http.StatusUnauthorized)
			}
			if respErr.Details != tt.details {
				t.Errorf("details = %q, want %q", respErr.Details, tt.details)
			}
		})
	}
}
