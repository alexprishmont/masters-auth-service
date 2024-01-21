package validation

import (
	"auth-sso/lib/validation/rules/documenttype"
	"auth-sso/lib/validation/rules/timestamp"
	"auth-sso/lib/validation/rules/uuid"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateStruct(someStruct interface{}) string {
	validate := validator.New()
	validate = registerCustomValidators(validate)

	var sb strings.Builder

	err := validate.Struct(someStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString("Field validation for '" + err.Field() + "' failed on the '" + err.Tag() + "' tag\n")
		}
	}

	return sb.String()
}

func registerCustomValidators(validator *validator.Validate) *validator.Validate {
	validator.RegisterValidation("uuid", uuid.Validate)
	validator.RegisterValidation("documenttype", documenttype.Validate)
	validator.RegisterValidation("timestamp", timestamp.Validate)

	return validator
}
