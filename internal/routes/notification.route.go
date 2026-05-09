package routes

import (
	"github.com/bigdann09/notifications/internal/controllers"
	"github.com/bigdann09/notifications/internal/services/notification"
	"github.com/gin-gonic/gin"
)

func (routes *Route) NotificationRoute(version *gin.RouterGroup) {
	routes.services.Logger.Info("notification controllers, services and routes")
	service := notification.NewNotificationService(routes.services.Database, routes.services.Logger)
	controller := controllers.NewNotificationController(service)

	notification := version.Group("/notifications")
	notification.GET("/", controller.FindAll)
}
