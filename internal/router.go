package internal

import (
	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/internal/health"
	authHttp "github.com/guvanch999/abrayly-mebel-backend/internal/internal/auth/http"
	bannerHttp "github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/http"
	categoryHttp "github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/http"
	productHttp "github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/http"
)

// InitRouters registers the routes of every module.
func InitRouters(server *cmd.Server) {
	server.Echo.GET("/health", health.Handler(server.Database.DB))

	authHttp.InitAuthRoutes(server)
	categoryHttp.InitCategoryRouter(server)
	productHttp.InitProductsRouter(server)
	bannerHttp.BannerRouter(server)
}
