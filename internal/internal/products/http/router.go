package http

import (
	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/handler"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/products/service"
	"github.com/guvanch999/abrayly-mebel-backend/internal/middleware"
)

func InitProductsRouter(server *cmd.Server) {
	productClientRouter := server.Echo.Group("/api/v1/products")
	productRepository := repository.NewProductRepository(server.Database)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	authAdminMiddleware := middleware.AuthAdminMiddleware(server)
	productAdminRouter := server.Echo.Group("/api/v1/admin/products", authAdminMiddleware)

	productAdminRouter.POST("", productHandler.CreateProduct)
	productAdminRouter.PUT("/:id", productHandler.UpdateProduct)
	productAdminRouter.DELETE("/:id", productHandler.DeleteProduct)

	productClientRouter.GET("", productHandler.GetProducts)
	productClientRouter.GET("/:id/recommended", productHandler.GetRecommendedProducts)
	productClientRouter.GET("/:id", productHandler.GetProduct)

	productAdminRouter.POST("/:id/images", productHandler.AddImageToProduct)
	productAdminRouter.DELETE("/:id/images/:image_id", productHandler.RemoveImageFromProduct)

	productAdminRouter.POST("/:id/attributes", productHandler.AddAttributeToProduct)
	productAdminRouter.DELETE("/:id/attributes/:attribute_id", productHandler.RemoveAttributeFromProduct)

	productAdminRouter.POST("/:id/attributes/:attribute_id/values/delete", productHandler.RemoveValueFromAttribute)
	productAdminRouter.POST("/:id/attributes/:attribute_id/values", productHandler.AddValueToAttribute)
}
