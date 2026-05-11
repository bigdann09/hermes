package apiresponse

import (
	"net/http"

	"github.com/bigdann09/notifications/pkgs/validators"
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int                     `json:"code"`
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Errors  []validators.FieldError `json:"errors,omitempty"`
	Data    any                     `json:"data,omitempty"`
}

func NewWithErrorFields(code int, message string, fields ...validators.FieldError) *APIResponse {
	return &APIResponse{
		Code:    code,
		Status:  statusText(code),
		Message: message,
		Errors:  fields,
	}
}

func NewWithData(code int, message string, data interface{}) *APIResponse {
	return &APIResponse{
		Code:    code,
		Status:  statusText(code),
		Message: message,
		Data:    data,
	}
}

func Response(c *gin.Context, response *APIResponse) {
	c.JSON(response.Code, response)
}

func Created(c *gin.Context, message string, data interface{}) {
	response := NewWithData(http.StatusCreated, message, data)
	c.JSON(http.StatusCreated, response)
}

func OK(c *gin.Context, message string, data interface{}) {
	response := NewWithData(http.StatusOK, message, data)
	c.JSON(http.StatusOK, response)
}

func InternalServerError(message string) *APIResponse {
	return NewWithErrorFields(http.StatusInternalServerError, message)
}

func BadRequest(message string) *APIResponse {
	return NewWithErrorFields(http.StatusBadRequest, message)
}

func NotFound(message string) *APIResponse {
	return NewWithErrorFields(http.StatusNotFound, message)
}

func Unauthorized(message string) *APIResponse {
	return NewWithErrorFields(http.StatusUnauthorized, message)
}

func Forbidden(message string) *APIResponse {
	return NewWithErrorFields(http.StatusForbidden, message)
}

func Conflict(message string) *APIResponse {
	return NewWithErrorFields(http.StatusConflict, message)
}

func Unprocessible(fields []validators.FieldError) *APIResponse {
	return NewWithErrorFields(http.StatusUnprocessableEntity, "validation failed", fields...)
}

func statusText(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return http.StatusText(code)
	}
}
