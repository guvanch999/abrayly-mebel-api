package middleware

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/admin/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
	"github.com/labstack/echo/v4"
)

// AdminContextKey is the echo.Context key holding the authenticated *models.Admin.
const AdminContextKey = "admin"

// AdminFinder loads an admin by ID.
type AdminFinder interface {
	FindByID(id uuid.UUID) (*models.Admin, error)
}

// AuthAdminMiddleware protects admin routes with an access token issued by
// POST /auth/admin/login.
func AuthAdminMiddleware(server *cmd.Server) echo.MiddlewareFunc {
	return RequireAdmin(server.Config.AdminJwtSecret, repository.NewAdminRepository(server.Database))
}

// RequireAdmin validates the "Authorization: Bearer <token>" header against
// secret and stores the matching admin in the request context.
func RequireAdmin(secret string, admins AdminFinder) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := BearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
			if err != nil {
				return exceptions.NewResponseError(exceptions.ErrUnauthorized, err)
			}

			adminID, err := util.CheckToken(token, secret)
			if err != nil {
				return err
			}

			admin, err := admins.FindByID(adminID)
			if err != nil {
				return exceptions.NewResponseError(exceptions.ErrUnauthorized, errors.New("admin not found"))
			}

			c.Set(AdminContextKey, admin)
			return next(c)
		}
	}
}

// BearerToken extracts the token from an Authorization header value.
func BearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("token is required")
	}
	scheme, token, found := strings.Cut(strings.TrimSpace(header), " ")
	token = strings.TrimSpace(token)
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", errors.New("invalid authorization header")
	}
	return token, nil
}
