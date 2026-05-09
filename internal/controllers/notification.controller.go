package controllers

import (
	"net/http"

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
	var query interface{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, query)
}
