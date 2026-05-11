package controllers

import (
	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/internal/services/notification"
	"github.com/bigdann09/notifications/pkgs/apiresponse"
	"github.com/bigdann09/notifications/pkgs/apiresponse/binder"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	service notification.INotificationService
}

func NewNotificationController(service notification.INotificationService) *NotificationController {
	return &NotificationController{service}
}

func (handler *NotificationController) FindAll(c *gin.Context) {
	var query dtos.NotificationQuery
	if err := binder.BindQuery(c, &query); err != nil {
		apiresponse.Response(c, err)
		return
	}

	result, err := handler.service.FindAll(query)
	if err != nil {
		apiresponse.Response(c, err)
		return
	}

	apiresponse.OK(c, "notifications retrieved successfully", result)
}
