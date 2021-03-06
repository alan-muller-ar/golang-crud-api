package validator

import (
	"fmt"

	"github.com/alan-muller-ar/alan-muller-ar-lahaus-backend/pkg/infrastructure/errors"
	enLocale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

const (
	MinValueValidationTag      = "min value violated"
	MaxValueValidationTag      = "max value violated"
	OutOfRangeValidationTAg    = "value out of range"
	InvalidStatusValidationTag = "invalid status"
)

type RequestValidator interface {
	ValidateCustom(v *validator.Validate, i interface{})
}

func Validate(o interface{}) error {
	en := enLocale.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")

	v := validator.New()

	return validate(v, translator, o)
}

func ValidateWithCustom(v *validator.Validate, o interface{}) error {
	en := enLocale.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")

	return validate(v, translator, o)
}

func validate(validate *validator.Validate, translator ut.Translator, o interface{}) error {
	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		return fmt.Errorf("error register default translation on validator. cause: %s", err.Error())
	}

	if err := translateOverride(validate, translator); err != nil {
		return fmt.Errorf("error register custom translation on validator. cause: %s", err.Error())
	}

	if err := validate.Struct(o); err != nil {
		var causes []string
		for _, e := range err.(validator.ValidationErrors) {
			causes = append(causes, e.Translate(translator))
		}

		return errors.NewMultiCauseError("invalid data", causes...)
	}

	return nil
}

func translateOverride(validate *validator.Validate, t ut.Translator) error {
	return validate.RegisterTranslation("required", t,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is required", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		},
	)
}
