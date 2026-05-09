package controllers

import (
	"github.com/bigdann09/notifications/internal/services/notification"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	service notification.INotificationService
}

func (handler *NotificationController) FindAll(c *gin.Context) {

}
