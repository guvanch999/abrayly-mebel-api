package http

import (
	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/admin/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/auth/handler"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/auth/service"
)

func InitAuthRoutes(server *cmd.Server) {
	authAdminRouter := server.Echo.Group("/auth/admin")

	adminRepository := repository.NewAdminRepository(server.Database)
	adminService := service.NewAdminService(adminRepository, server.Config)
	adminHandler := handler.NewAuthHandler(adminService)

	authAdminRouter.POST("/login", adminHandler.Login)
	authAdminRouter.POST("/refresh", adminHandler.RefreshToken)
}
