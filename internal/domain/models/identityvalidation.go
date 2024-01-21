package models

import (
	identityverificationv1 "github.com/alexprishmont/masters-protos/gen/go/identityverification"
	"time"
)

type IdentityValidation struct {
	ValidationId string                              `bson:"validationId"`
	User         User                                `bson:"user"`
	DocumentType identityverificationv1.DocumentType `bson:"documentType"`
	Status       identityverificationv1.Status       `bson:"status"`
	CreatedAt    time.Time                           `bson:"createdAt"`
	UpdatedAt    time.Time                           `bson:"updatedAt"`
}
