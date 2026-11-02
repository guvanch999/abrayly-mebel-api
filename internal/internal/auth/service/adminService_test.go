package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/config"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/admin/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
)

// fakeAdmins is an in-memory AdminStore.
type fakeAdmins struct {
	admins []*models.Admin
	err    error
}

func (f *fakeAdmins) GetAdminByUsername(username string) (*models.Admin, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, a := range f.admins {
		if a.Username == username {
			return a, nil
		}
	}
	return nil, repository.ErrAdminNotFound
}

func (f *fakeAdmins) FindByID(id uuid.UUID) (*models.Admin, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, a := range f.admins {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, repository.ErrAdminNotFound
}

var testConfig = &config.Schema{
	AdminJwtSecret:     "access-secret-for-tests",
	AdminRefreshSecret: "refresh-secret-for-tests",
}

func newService(t *testing.T) (*AdminService, *models.Admin) {
	t.Helper()
	hash, err := util.HashString("correct-password")
	if err != nil {
		t.Fatalf("HashString: %v", err)
	}
	admin := &models.Admin{ID: uuid.New(), Name: "Owner", Username: "owner", Password: hash}
	return NewAdminService(&fakeAdmins{admins: []*models.Admin{admin}}, testConfig), admin
}

func assertStatus(t *testing.T, err error, status int) {
	t.Helper()
	var respErr *exceptions.ResponseError
	if !errors.As(err, &respErr) {
		t.Fatalf("error = %v (%T), want *exceptions.ResponseError", err, err)
	}
	if respErr.Code != status {
		t.Fatalf("status = %d, want %d", respErr.Code, status)
	}
}

func TestLoginIssuesSeparateAccessAndRefreshTokens(t *testing.T) {
	svc, admin := newService(t)

	tokens, err := svc.Login("owner", "correct-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if id, err := util.CheckToken(tokens.AccessToken, testConfig.AdminJwtSecret); err != nil || id != admin.ID {
		t.Errorf("access token: id=%s err=%v, want %s", id, err, admin.ID)
	}
	if id, err := util.CheckToken(tokens.RefreshToken, testConfig.AdminRefreshSecret); err != nil || id != admin.ID {
		t.Errorf("refresh token: id=%s err=%v, want %s", id, err, admin.ID)
	}

	// Each token is only valid with its own secret.
	if _, err := util.CheckToken(tokens.RefreshToken, testConfig.AdminJwtSecret); err == nil {
		t.Error("refresh token was accepted as an access token")
	}
	if _, err := util.CheckToken(tokens.AccessToken, testConfig.AdminRefreshSecret); err == nil {
		t.Error("access token was accepted as a refresh token")
	}
}

func TestLoginRejectsBadCredentialsUniformly(t *testing.T) {
	svc, _ := newService(t)

	_, wrongPassword := svc.Login("owner", "wrong-password")
	assertStatus(t, wrongPassword, http.StatusUnauthorized)

	_, unknownUser := svc.Login("nobody", "correct-password")
	assertStatus(t, unknownUser, http.StatusUnauthorized)

	if wrongPassword.Error() != unknownUser.Error() {
		t.Errorf("responses differ and reveal which usernames exist:\n%v\n%v", wrongPassword, unknownUser)
	}
}

func TestLoginPassesThroughStorageErrors(t *testing.T) {
	dbDown := errors.New("database is down")
	svc := NewAdminService(&fakeAdmins{err: dbDown}, testConfig)

	if _, err := svc.Login("owner", "correct-password"); !errors.Is(err, dbDown) {
		t.Fatalf("error = %v, want %v", err, dbDown)
	}
}

func TestRefreshToken(t *testing.T) {
	svc, admin := newService(t)

	tokens, err := svc.Login("owner", "correct-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	refreshed, err := svc.RefreshToken(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if id, err := util.CheckToken(refreshed.AccessToken, testConfig.AdminJwtSecret); err != nil || id != admin.ID {
		t.Errorf("refreshed access token: id=%s err=%v", id, err)
	}

	_, err = svc.RefreshToken(tokens.AccessToken)
	assertStatus(t, err, http.StatusUnauthorized)
}

func TestRefreshTokenForDeletedAdmin(t *testing.T) {
	svc, _ := newService(t)

	token, err := util.IssueToken(uuid.New(), testConfig.AdminRefreshSecret, RefreshTokenTTL)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	_, err = svc.RefreshToken(token)
	assertStatus(t, err, http.StatusUnauthorized)
}
