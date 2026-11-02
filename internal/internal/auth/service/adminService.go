package service

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/config"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/admin/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/auth/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
)

const (
	// AccessTokenTTL is the lifetime of an admin access token.
	AccessTokenTTL = 30 * time.Minute
	// RefreshTokenTTL is the lifetime of an admin refresh token.
	RefreshTokenTTL = 24 * time.Hour
)

// AdminStore is the persistence the auth service depends on.
type AdminStore interface {
	GetAdminByUsername(username string) (*models.Admin, error)
	FindByID(id uuid.UUID) (*models.Admin, error)
}

type AdminService struct {
	adminRepository AdminStore
	cfg             *config.Schema
}

func NewAdminService(adminRepository AdminStore, cfg *config.Schema) *AdminService {
	return &AdminService{adminRepository: adminRepository, cfg: cfg}
}

var errInvalidCredentials = errors.New("invalid username or password")

// Login checks the credentials and returns a fresh access/refresh token pair.
// Unknown usernames and wrong passwords produce the same 401 response.
func (s *AdminService) Login(username string, password string) (*dto.TokenResponseDto, error) {
	admin, err := s.adminRepository.GetAdminByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrAdminNotFound) {
			// Spend comparable time to a real check so response timing does
			// not reveal which usernames exist.
			_ = util.CompareString(dummyHash(), password)
			return nil, exceptions.NewResponseError(exceptions.ErrUnauthorized, errInvalidCredentials)
		}
		return nil, err
	}

	if err := util.CompareString(admin.Password, password); err != nil {
		return nil, exceptions.NewResponseError(exceptions.ErrUnauthorized, errInvalidCredentials)
	}

	return s.GenerateToken(admin)
}

// GenerateToken issues an access token signed with ADMIN_JWT_SECRET and a
// refresh token signed with ADMIN_REFRESH_SECRET, so neither can be used in
// place of the other.
func (s *AdminService) GenerateToken(admin *models.Admin) (*dto.TokenResponseDto, error) {
	accessToken, err := util.IssueToken(admin.ID, s.cfg.AdminJwtSecret, AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.IssueToken(admin.ID, s.cfg.AdminRefreshSecret, RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken exchanges a valid refresh token for a new token pair.
func (s *AdminService) RefreshToken(refreshToken string) (*dto.TokenResponseDto, error) {
	adminID, err := util.CheckToken(refreshToken, s.cfg.AdminRefreshSecret)
	if err != nil {
		return nil, err
	}

	admin, err := s.adminRepository.FindByID(adminID)
	if err != nil {
		if errors.Is(err, repository.ErrAdminNotFound) {
			return nil, exceptions.NewResponseError(exceptions.ErrUnauthorized, errors.New("admin not found"))
		}
		return nil, err
	}

	return s.GenerateToken(admin)
}

var (
	dummyHashOnce  sync.Once
	dummyHashValue string
)

func dummyHash() string {
	dummyHashOnce.Do(func() {
		dummyHashValue, _ = util.HashString(uuid.NewString())
	})
	return dummyHashValue
}
