package routes

import (
	"github.com/bigdann09/notifications/internal/controllers"
	"github.com/gin-gonic/gin"
)

func (routes *Route) NotificationRoute(version *gin.RouterGroup) {
	routes.services.Logger.Info("notification routes")
	controller := controllers.NotificationController{}

	notification := version.Group("/notifications")
	notification.GET("/", controller.FindAll)
}
