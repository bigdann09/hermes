package validators

import (
	"github.com/bigdann09/notifications/internal/models"
	"github.com/go-playground/validator/v10"
)

func (v *Validator) HasNotificationType(fl validator.FieldLevel) bool {
	if fl.Field().Len() == 0 {
		return true
	}

	value, ok := fl.Field().Interface().(models.NotificationType)
	if !ok {
		return false
	}

	switch value {
	case models.Social, models.System:
		return true
	default:
		return false
	}
}
