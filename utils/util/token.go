package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
)

// IssueToken signs an HS256 JWT whose subject is the given ID and which
// expires after ttl.
func IssueToken(subject uuid.UUID, secret string, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("token secret is empty")
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   subject.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// CheckToken verifies an HS256 JWT signed with secret and returns its subject.
// Tokens with another algorithm, a missing or past expiry, or a subject that
// is not a UUID are rejected with a 401 ResponseError.
func CheckToken(token, secret string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, unauthorized("token is required")
	}

	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		claims,
		func(*jwt.Token) (interface{}, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, unauthorized("token expired")
		}
		return uuid.Nil, unauthorized("invalid token")
	}

	subject, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, unauthorized("invalid token")
	}
	return subject, nil
}

func unauthorized(details string) error {
	return exceptions.NewResponseError(exceptions.ErrUnauthorized, errors.New(details))
}
