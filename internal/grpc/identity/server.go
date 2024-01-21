package identitygrpc

import (
	"auth-sso/lib/validation"
	"context"
	"encoding/json"
	identityverificationv1 "github.com/alexprishmont/masters-protos/gen/go/identityverification"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type DocumentType string

type Verification interface {
	StartValidation(ctx context.Context,
		userId string,
		documentType string,
	) (response ValidationResponse, err error)
	Status(ctx context.Context,
		validationId string,
	) (response StatusResponse, err error)
	DocumentUpload(ctx context.Context,
		validationId string,
		document []byte,
		documentFormat string,
	) (response DocumentUploadResponse, err error)
	EndValidation(ctx context.Context,
		validationId string,
	) (response EndValidationResponse, err error)
	UpdateValidation(ctx context.Context,
		validationId string,
		updatedInformation *UpdatedInfo,
	) (response UpdateValidationResponse, err error)
	CancelValidation(ctx context.Context,
		validationId string,
	) (response CancelValidationResponse, err error)
}

type ValidationRequest struct {
	UserId       string                              `validate:"required,uuid"`
	DocumentType identityverificationv1.DocumentType `validate:"required"`
}

type ValidationResponse struct {
	ValidationId string
	Status       identityverificationv1.Status
	Message      string
}

type StatusRequest struct {
	ValidationId string `validate:"required,uuid"`
}

type StatusResponse struct {
	Status      identityverificationv1.Status
	LastUpdated *timestamp.Timestamp
	Message     string
}

type DocumentUploadRequest struct {
	ValidationId   string `validate:"required,uuid"`
	Document       []byte `validate:"required"`
	DocumentFormat string `validate:"required"`
}

type DocumentUploadResponse struct {
	UploadStatus identityverificationv1.Status
	Message      string
}

type EndValidationRequest struct {
	ValidationId string `validate:"required,uuid"`
}

type EndValidationResponse struct {
	FinalStatus identityverificationv1.Status
	Message     string
}

type UpdateValidationRequest struct {
	ValidationId       string `validate:"required,uuid"`
	UpdatedInformation string `validate:"required"`
}

type UpdateValidationResponse struct {
	UpdateStatus identityverificationv1.Status
	Message      string
}

type CancelValidationRequest struct {
	ValidationId string `validate:"required,uuid"`
}

type CancelValidationResponse struct {
	CancellationStatus identityverificationv1.Status
	Message            string
}

type UpdatedInfo struct {
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
}

type serverAPI struct {
	identityverificationv1.UnimplementedIdentityValidationServer
	log                  *slog.Logger
	identityVerification Verification
}

func Register(gRPC *grpc.Server, log *slog.Logger, identityVerification Verification) {
	identityverificationv1.RegisterIdentityValidationServer(gRPC, &serverAPI{
		log:                  log,
		identityVerification: identityVerification,
	})
}

func (s *serverAPI) StartValidation(
	ctx context.Context,
	request *identityverificationv1.ValidationRequest,
) (*identityverificationv1.ValidationResponse, error) {
	req := ValidationRequest{
		UserId:       request.GetUserId(),
		DocumentType: request.GetDocumentType(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	response, err := s.identityVerification.StartValidation(ctx, req.UserId, req.DocumentType.String())

	if err != nil {
		// TODO: Handle error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityverificationv1.ValidationResponse{
		ValidationId: response.ValidationId,
		Status:       response.Status.String(),
		Message:      response.Message,
	}, nil
}

func (s *serverAPI) Status(
	ctx context.Context,
	request *identityverificationv1.StatusRequest,
) (*identityverificationv1.StatusResponse, error) {
	req := StatusRequest{
		ValidationId: request.GetValidationId(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	response, err := s.identityVerification.Status(ctx, req.ValidationId)

	if err != nil {
		// TODO: Handle error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityverificationv1.StatusResponse{
		Status:      response.Status.String(),
		LastUpdated: response.LastUpdated,
		Message:     response.Message,
	}, nil
}

func (s *serverAPI) DocumentUpload(
	ctx context.Context,
	request *identityverificationv1.DocumentUploadRequest,
) (*identityverificationv1.DocumentUploadResponse, error) {
	req := DocumentUploadRequest{
		ValidationId:   request.GetValidationId(),
		Document:       request.GetDocument(),
		DocumentFormat: request.GetDocumentFormat(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	response, err := s.identityVerification.DocumentUpload(ctx, req.ValidationId, req.Document, req.DocumentFormat)

	if err != nil {
		// TODO: Handle error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityverificationv1.DocumentUploadResponse{
		UploadStatus: response.UploadStatus.String(),
		Message:      response.Message,
	}, nil
}

func (s *serverAPI) EndValidation(
	ctx context.Context,
	request *identityverificationv1.EndValidationRequest,
) (*identityverificationv1.EndValidationResponse, error) {
	req := EndValidationRequest{
		ValidationId: request.GetValidationId(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	response, err := s.identityVerification.EndValidation(ctx, req.ValidationId)

	if err != nil {
		// TODO: Handle error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityverificationv1.EndValidationResponse{
		FinalStatus: response.FinalStatus.String(),
		Message:     response.Message,
	}, nil
}

func (s *serverAPI) UpdateValidation(
	ctx context.Context,
	request *identityverificationv1.UpdateValidationRequest,
) (*identityverificationv1.UpdateValidationResponse, error) {
	req := UpdateValidationRequest{
		ValidationId:       request.GetValidationId(),
		UpdatedInformation: request.GetUpdatedInformation(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	parsedUpdatedInfo, err := parseUpdatedInformation(req.UpdatedInformation)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	response, err := s.identityVerification.UpdateValidation(ctx, req.ValidationId, parsedUpdatedInfo)

	if err != nil {
		// TODO: Handle error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityverificationv1.UpdateValidationResponse{
		UpdateStatus: response.UpdateStatus.String(),
		Message:      response.Message,
	}, nil
}

func (s *serverAPI) CancelValidation(
	ctx context.Context,
	request *identityverificationv1.CancelValidationRequest,
) (*identityverificationv1.CancelValidationResponse, error) {
	req := CancelValidationRequest{
		ValidationId: request.GetValidationId(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	response, err := s.identityVerification.CancelValidation(ctx, req.ValidationId)

	if err != nil {
		// TODO: Handle error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityverificationv1.CancelValidationResponse{
		CancellationStatus: response.CancellationStatus.String(),
		Message:            response.Message,
	}, nil
}

func parseUpdatedInformation(jsonString string) (*UpdatedInfo, error) {
	var info UpdatedInfo
	err := json.Unmarshal([]byte(jsonString), &info)

	if err != nil {
		return nil, err
	}

	return &info, nil
}
