package http

import (
	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/handler"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/category/service"
	"github.com/guvanch999/abrayly-mebel-backend/internal/middleware"
)

func InitCategoryRouter(server *cmd.Server) {
	categoryRouter := server.Echo.Group("/api/v1/categories")

	categoryRepository := repository.NewCategoryRepository(server.Database)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	adminMiddleware := middleware.AuthAdminMiddleware(server)

	categoryRouter.GET("/:id", categoryHandler.GetCategory)
	categoryRouter.GET("", categoryHandler.GetAllCategories)

	adminRouter := server.Echo.Group("/api/v1/admin/categories", adminMiddleware)

	adminRouter.POST("", categoryHandler.CreateCategory)
	adminRouter.PUT("/:id", categoryHandler.UpdateCategory)
	adminRouter.DELETE("/:id", categoryHandler.DeleteCategory)
}
