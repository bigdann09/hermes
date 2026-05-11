package validators

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateError(err error) []FieldError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}

	translated := ve.Translate(Translator)
	out := make([]FieldError, 0, len(translated))
	for field, message := range translated {
		out = append(out, FieldError{
			Field:   field,
			Message: message,
		})
	}
	return out
}
