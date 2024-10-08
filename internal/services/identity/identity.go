package identity

import (
	"auth-sso/internal/domain/models"
	identitygrpc "auth-sso/internal/grpc/identity"
	"auth-sso/internal/storage"
	"auth-sso/internal/tasks/handlers/identity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	identityverificationv1 "github.com/alexprishmont/masters-protos/gen/go/identityverification"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type Verification struct {
	log                *slog.Logger
	asynqClient        *asynq.Client
	userProvider       UserProvider
	validationSaver    ValidationSaver
	validationProvider ValidationProvider
}

type ValidationSaver interface {
	CreateNewValidation(ctx context.Context,
		user models.User,
		documentType string,
	) (string, error)
}

type UserProvider interface {
	UserById(ctx context.Context, id string) (models.User, error)
}

type ValidationProvider interface {
	Validation(ctx context.Context, id string) (models.IdentityValidation, error)
	DoesValidationExist(ctx context.Context, userId string) (bool, error)
}

var (
	ErrorInvalidUserId = errors.New("invalid user id")
)

func New(
	log *slog.Logger,
	asynqClient *asynq.Client,
	userProvider UserProvider,
	validationSaver ValidationSaver,
	validationProvider ValidationProvider,
) *Verification {
	return &Verification{
		log:                log,
		asynqClient:        asynqClient,
		userProvider:       userProvider,
		validationSaver:    validationSaver,
		validationProvider: validationProvider,
	}
}

// StartValidation initiates a validation processes as a async process.
// Creates new record in database
// Dispatches async job which performs user information checks
func (v *Verification) StartValidation(
	ctx context.Context,
	userId string,
	documentType string,
) (identitygrpc.ValidationResponse, error) {
	const op = "identity.startValidation"

	log := v.log.With(
		slog.String("op", op),
	)

	log.Info("Starting validation process", slog.String("userId", userId))

	user, err := v.userProvider.UserById(ctx, userId)

	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			v.log.Warn("user not found", slog.String("error", err.Error()))

			return identitygrpc.ValidationResponse{}, fmt.Errorf("%s: %w", op, ErrorInvalidUserId)
		}

		v.log.Error("failed to get user", slog.String("error", err.Error()))

		return identitygrpc.ValidationResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	exists, err := v.validationProvider.DoesValidationExist(ctx, user.UniqueId)

	if exists && err == nil {
		v.log.Error("validation already exists for the user", slog.String("user", user.UniqueId))

		return identitygrpc.ValidationResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := v.validationSaver.CreateNewValidation(ctx, user, documentType)

	if err != nil {
		log.Error("Failed to create new validation process.", slog.String("error", err.Error()))

		return identitygrpc.ValidationResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	payload := identity.VerificationTaskPayload{
		User: user,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("Failed to marshal payload", slog.String("userId", userId))

		return identitygrpc.ValidationResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	task := asynq.NewTask(identity.TaskIdentifier, payloadBytes)
	if _, err := v.asynqClient.Enqueue(task); err != nil {
		log.Error("Failed to dispatch identity validation task.", slog.String("userId", userId))

		return identitygrpc.ValidationResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Identity verification process is started.", slog.String("userId", userId))

	return identitygrpc.ValidationResponse{
		ValidationId: id,
		Status:       identityverificationv1.Status_PENDING,
		Message:      "Identity validation process is started.",
	}, nil
}

func (v *Verification) Status(
	ctx context.Context,
	validationId string,
) (identitygrpc.StatusResponse, error) {
	const op = "identity.Status"

	log := v.log.With(
		slog.String("op", op),
	)

	log.Info("receiving identity validation status")

	validation, err := v.validationProvider.Validation(ctx, validationId)

	if err != nil {
		log.Error("Failed to retrieve validation object.", slog.String("error", err.Error()))

		return identitygrpc.StatusResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return identitygrpc.StatusResponse{
		Status:      validation.Status,
		LastUpdated: timestamppb.New(validation.UpdatedAt),
		Message:     "OK, validation status.",
	}, nil
}

func (v *Verification) DocumentUpload(
	ctx context.Context,
	validationId string,
	document []byte,
	documentFormat string,
) (identitygrpc.DocumentUploadResponse, error) {

	return identitygrpc.DocumentUploadResponse{}, nil
}

func (v *Verification) EndValidation(
	ctx context.Context,
	validationId string,
) (identitygrpc.EndValidationResponse, error) {

	return identitygrpc.EndValidationResponse{}, nil
}

func (v *Verification) UpdateValidation(
	ctx context.Context,
	validationId string,
	updatedInformation *identitygrpc.UpdatedInfo,
) (identitygrpc.UpdateValidationResponse, error) {

	return identitygrpc.UpdateValidationResponse{}, nil
}

func (v *Verification) CancelValidation(
	ctx context.Context,
	validationId string,
) (identitygrpc.CancelValidationResponse, error) {

	return identitygrpc.CancelValidationResponse{}, nil
}
