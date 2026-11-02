package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Response is the body returned by the health endpoint.
type Response struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// DBProvider returns the underlying connection pool (gorm.DB.DB satisfies it).
type DBProvider func() (*sql.DB, error)

// Handler reports whether the service and its database are reachable.
//
// @Summary Health check
// @Description Liveness/readiness probe; pings PostgreSQL
// @Tags system
// @Produce json
// @Success 200 {object} health.Response
// @Failure 503 {object} health.Response
// @Router /health [get]
func Handler(db DBProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db()
		if err == nil {
			err = sqlDB.PingContext(ctx)
		}
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusServiceUnavailable, Response{Status: "unavailable", Database: "down"})
		}
		return c.JSON(http.StatusOK, Response{Status: "ok", Database: "up"})
	}
}
