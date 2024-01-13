package identity

import (
	identitygrpc "auth-sso/internal/grpc/identity"
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log/slog"
)

type Verification struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) *Verification {
	return &Verification{
		log: log,
	}
}

func (v *Verification) StartValidation(
	ctx context.Context,
	userId string,
	documentType string,
	initiationTimestamp *timestamp.Timestamp,
) (identitygrpc.ValidationResponse, error) {

	return identitygrpc.ValidationResponse{}, nil
}

func (v *Verification) Status(
	ctx context.Context, validationId string) (identitygrpc.StatusResponse, error) {
	panic("implement me")
}

func (v *Verification) DocumentUpload(
	ctx context.Context,
	validationId string,
	document []byte,
	documentFormat string,
) (identitygrpc.DocumentUploadResponse, error) {
	panic("implement me")
}

func (v *Verification) EndValidation(
	ctx context.Context,
	validationId string,
) (identitygrpc.EndValidationResponse, error) {
	panic("implement me")
}

func (v *Verification) UpdateValidation(
	ctx context.Context,
	validationId string,
	updatedInformation *identitygrpc.UpdatedInfo,
) (identitygrpc.UpdateValidationResponse, error) {
	panic("implement me")
}

func (v *Verification) CancelValidation(
	ctx context.Context,
	validationId string,
) (identitygrpc.CancelValidationResponse, error) {
	panic("implement me")
}
