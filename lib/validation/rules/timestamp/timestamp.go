package timestamp

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampValidator(fl validator.FieldLevel) bool {
	ts, ok := fl.Field().Interface().(*timestamppb.Timestamp)
	if !ok || ts == nil {
		return false
	}

	return true
}
