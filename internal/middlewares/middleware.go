package middlewares

import (
	"strings"
	"time"

	"github.com/bigdann09/notifications/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Middleware struct {
	logger *zap.Logger
	config *config.Config
}

func NewMiddleware(logger *zap.Logger, cfg *config.Config) *Middleware {
	return &Middleware{
		logger: logger,
		config: cfg,
	}
}

func (middleware *Middleware) CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(middleware.config.App.AllowedOrigins, ","),
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
