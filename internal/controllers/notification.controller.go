package controllers

import (
	"github.com/bigdann09/notifications/internal/services/notification"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	service notification.INotificationService
}

func NewNotificationController(service notification.INotificationService) *NotificationController {
	return &NotificationController{service}
}

func (handler *NotificationController) FindAll(c *gin.Context) {

}
