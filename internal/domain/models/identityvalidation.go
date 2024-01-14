package models

import (
	identityverificationv1 "github.com/alexprishmont/masters-protos/gen/go/identityverification"
	"time"
)

type IdentityValidation struct {
	ValidationId string
	User         User
	DocumentType identityverificationv1.DocumentType
	Status       identityverificationv1.Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
