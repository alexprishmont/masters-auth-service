package validation

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateStruct(someStruct interface{}) string {
	validate := validator.New()
	var sb strings.Builder

	err := validate.Struct(someStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString("Field validation for '" + err.Field() + "' failed on the '" + err.Tag() + "' tag\n")
		}
	}

	return sb.String()
}
