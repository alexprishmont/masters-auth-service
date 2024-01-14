package models

type Document struct {
	Validation     IdentityValidation
	Document       []byte
	DocumentFormat string
}
