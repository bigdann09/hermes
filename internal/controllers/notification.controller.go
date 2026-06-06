package controllers

import (
	"fmt"

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

func (handler *NotificationController) Create(c *gin.Context) {
	var request dtos.NotificationRequest
	if err := binder.BindJSON(c, &request); err != nil {
		apiresponse.Response(c, err)
		return
	}

	fmt.Println("request", request)

	response := handler.service.Send(request)
	apiresponse.Response(c, response)
}
