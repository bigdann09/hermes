package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

type ValidationField struct {
	Tag  string
	Func validator.Func
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Register(fields ...*ValidationField) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, field := range fields {
			v.RegisterValidation(field.Tag, field.Func)
		}
	}
}
