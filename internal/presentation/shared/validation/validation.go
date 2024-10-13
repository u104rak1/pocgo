package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/pocgo/pkg/email"
)

func SetupCustomValidation(e *echo.Echo) {
	v := validator.New()
	v.RegisterValidation("validEmail", validEmail)
	e.Validator = &CustomValidator{Validator: v}
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func validEmail(fl validator.FieldLevel) bool {
	return email.IsValid(fl.Field().String())
}
