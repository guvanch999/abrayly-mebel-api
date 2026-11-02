package middleware

import (
	"net/http"

	"github.com/guvanch999/abrayly-mebel-backend/cmd"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// RegisterMiddlewares installs the global middleware chain.
func RegisterMiddlewares(server *cmd.Server) {
	server.Echo.Pre(echoMiddleware.RemoveTrailingSlash())
	server.Echo.Use(echoMiddleware.Recover())
	server.Echo.Use(
		echoMiddleware.BodyLimitWithConfig(
			echoMiddleware.BodyLimitConfig{
				Limit: "250M",
			},
		),
	)

	server.Echo.Use(
		echoMiddleware.CORSWithConfig(
			echoMiddleware.CORSConfig{
				AllowOrigins: server.Config.AllowedOrigins(),
				AllowMethods: []string{
					http.MethodGet,
					http.MethodPut,
					http.MethodPatch,
					http.MethodPost,
					http.MethodDelete,
					http.MethodOptions,
				},
				AllowHeaders: []string{
					echo.HeaderOrigin,
					echo.HeaderContentType,
					echo.HeaderAccept,
					echo.HeaderAuthorization,
				},
			},
		),
	)
	server.Echo.Use(
		echoMiddleware.LoggerWithConfig(
			echoMiddleware.LoggerConfig{
				Format:           "${time_custom} | ${method} ${uri} | Status: ${status} | Error: ${error} \n",
				CustomTimeFormat: "2006/01/02 15:04:05",
			},
		),
	)
	// 20 requests per second per client IP, kept in memory.
	server.Echo.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
}
