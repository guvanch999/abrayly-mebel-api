package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString returns the bcrypt hash of str.
func HashString(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CompareString returns nil when plain matches the bcrypt hash.
func CompareString(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
