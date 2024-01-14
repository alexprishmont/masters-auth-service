package documenttype

import "github.com/go-playground/validator/v10"

var allowedDocumentTypes = map[string]bool{
	"PASSPORT":       true,
	"DRIVER_LICENSE": true,
	"NATIONAL_ID":    true,
}

func TypeValidation(fl validator.FieldLevel) bool {
	docType := fl.Field().String()
	_, exists := allowedDocumentTypes[docType]

	return exists
}
