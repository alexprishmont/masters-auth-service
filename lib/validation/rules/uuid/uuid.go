package uuid

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Validate(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())

	return err == nil
}
