package module

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/middleware"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func registerRoutes(r *gin.RouterGroup, config *config.Config) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}

var RouterModule = fx.Options(
	fx.Provide(func() *gin.Engine {
		r := gin.Default()

		r.Use(cors.Default())
		r.Use(middleware.ErrorHandler())

		return r
	}),
	fx.Provide(func(g *gin.Engine) *gin.RouterGroup {
		return g.Group("/api/v1")
	}),

	fx.Invoke(registerRoutes),
)
