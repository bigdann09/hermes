package routes

import (
	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/middlewares"
	"github.com/bigdann09/notifications/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Route struct {
	engine     *gin.Engine
	middleware *middlewares.Middleware
	services   *services.Service
}

func NewRoute(services *services.Service, engine *gin.Engine, logger *zap.Logger, cfg *config.Config) *Route {
	return &Route{
		services:   services,
		engine:     engine,
		middleware: middlewares.NewMiddleware(logger, cfg),
	}
}

func (routes *Route) Register() *gin.Engine {
	routes.engine.Use(routes.middleware.CORS())

	version := routes.engine.Group("/api/v1")
	version.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	routes.NotificationRoute(version)
	return routes.engine
}
