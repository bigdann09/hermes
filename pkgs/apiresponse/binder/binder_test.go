package binder_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/pkgs/api_response/binder"
	"github.com/bigdann09/notifications/pkgs/validators"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBindQuery_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Initialize validator
	v := validators.NewValidator()
	v.Register(
		&validators.ValidationField{
			Tag:         "has_notification_type",
			Func:        v.HasNotificationType,
			Translation: "{0} is not a supported notification type",
		},
	)

	t.Run("invalid type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/?type=invalid", nil)

		var query dtos.NotificationQuery
		resp := binder.BindQuery(c, &query)

		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Len(t, resp.Errors, 1)
		assert.Equal(t, "Type", resp.Errors[0].Field)
	})

	t.Run("valid type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/?type=system", nil)

		var query dtos.NotificationQuery
		resp := binder.BindQuery(c, &query)

		assert.Nil(t, resp)
		assert.Equal(t, "system", string(*query.Type))
	})
}
