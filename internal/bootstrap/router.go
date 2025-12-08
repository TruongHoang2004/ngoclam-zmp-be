package bootstrap

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/controllers"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func registerRoutes(
	r *gin.RouterGroup,
	imageController *controllers.ImageController,
	folderController *controllers.FolderController,
	productController *controllers.ProductController,
	categoryController *controllers.CategoryController,
	authController *controllers.AuthController,
) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	imageController.RegisterRoutes(r)
	folderController.RegisterRoutes(r)
	productController.RegisterRoutes(r)
	categoryController.RegisterRoutes(r)
	authController.RegisterRoutes(r)
}

var RouterModule = fx.Options(
	fx.Provide(func() *gin.Engine {
		r := gin.Default()

		r.Use(cors.Default())

		return r
	}),
	fx.Provide(func(g *gin.Engine) *gin.RouterGroup {
		return g.Group("/api/v1")
	}),

	fx.Invoke(registerRoutes),
)
