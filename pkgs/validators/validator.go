package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Translator ut.Translator

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

type ValidationField struct {
	Tag         string
	Func        validator.Func
	Translation string
}

func NewValidator() *Validator {
	validator, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	locale := en.New()
	uni_locale := ut.New(locale, locale)
	translator, _ := uni_locale.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validator, translator)

	Translator = translator
	return &Validator{
		validator:  validator,
		translator: translator,
	}
}

func (v *Validator) Register(fields ...*ValidationField) {
	for _, field := range fields {
		v.validator.RegisterValidation(field.Tag, field.Func)

		v.validator.RegisterTranslation(field.Tag, v.translator,
			func(ut ut.Translator) error {
				return ut.Add(field.Tag, field.Translation, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(field.Tag, fe.Field())
				return t
			},
		)
	}
}
