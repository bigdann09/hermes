package binder

import (
	"github.com/bigdann09/notifications/pkgs/apiresponse"
	"github.com/bigdann09/notifications/pkgs/validators"
	"github.com/gin-gonic/gin"
)

func BindQuery(c *gin.Context, payload any) *apiresponse.APIResponse {
	if err := c.ShouldBindQuery(payload); err != nil {
		fields := validators.ValidateError(err)
		return apiresponse.Unprocessible(fields)
	}
	return nil
}

func BindJSON(c *gin.Context, payload any) *apiresponse.APIResponse {
	if err := c.ShouldBindJSON(payload); err != nil {
		fields := validators.ValidateError(err)
		return apiresponse.Unprocessible(fields)
	}
	return nil
}

func BindParam(c *gin.Context, payload any) *apiresponse.APIResponse {
	if err := c.ShouldBindUri(payload); err != nil {
		fields := validators.ValidateError(err)
		return apiresponse.Unprocessible(fields)
	}
	return nil
}
