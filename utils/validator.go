package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func ValidateStruct(s interface{}) map[string]string {
	err := Validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errors[field] = getErrorMessage(e)
	}
	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return "minimum length is " + e.Param()
	case "max":
		return "maximum length is " + e.Param()
	case "gt":
		return "must be greater than " + e.Param()
	case "uuid":
		return "invalid UUID format"
	case "alphanum":
		return "only letters and numbers allowed"
	case "nefield":
		return "must be different from " + e.Param()
	default:
		return "invalid value"
	}
}