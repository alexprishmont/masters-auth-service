package documenttype

import "github.com/go-playground/validator/v10"

var allowedDocumentTypes = map[string]bool{
	"passport":       true,
	"driver_license": true,
	"national_id":    true,
}

func DocumentTypeValidation(fl validator.FieldLevel) bool {
	docType := fl.Field().String()
	_, exists := allowedDocumentTypes[docType]

	return exists
}
