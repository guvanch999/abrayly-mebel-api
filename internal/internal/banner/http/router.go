package http

import (
	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/handler"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/service"
	"github.com/guvanch999/abrayly-mebel-backend/internal/middleware"
)

func BannerRouter(server *cmd.Server) {
	bannerRouter := server.Echo.Group("/banners")

	bannerRepository := repository.NewBannerRepository(server.Database)
	bannerService := service.NewBannerService(bannerRepository)
	bannerHandler := handler.NewBannerHandler(bannerService)

	bannerRouter.GET("", bannerHandler.GetAll)
	bannerRouter.GET("/:id", bannerHandler.GetByID)

	bannerAdminGroup := server.Echo.Group("/admin/banners")
	adminAuthMiddleware := middleware.AuthAdminMiddleware(server)
	bannerAdminGroup.Use(adminAuthMiddleware)

	bannerAdminGroup.POST("", bannerHandler.Create)
	bannerAdminGroup.PUT("/:id", bannerHandler.Update)
	bannerAdminGroup.DELETE("/:id", bannerHandler.Delete)
}
