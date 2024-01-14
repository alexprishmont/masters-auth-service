package identity

import (
	"auth-sso/internal/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

type VerificationTaskPayload struct {
	User models.User
}

const TaskIdentifier = "identity:validate"

func HandleIdentityVerificationTask(ctx context.Context, task *asynq.Task) error {
	const op = "tasks.handlers.identity.HandleIdentityVerificationTask"

	var payload VerificationTaskPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Call user management component to retrieve user data
	// TODO: Validate received user data on presence of:
	// TODO: Full name, date of birth, gender, email address, address, phone number, national identification code
	// TODO: Validate person photo
	// TODO: If person photo not found in database require it from the person by changing verification status and sending email
	// TODO: If person photo is found do validation with the photo on the provided document
	// TODO: If document is not provided ask person to do that
	// TODO: Check person's photo on deepfake presence using ML component

	fmt.Println("Processed job...")

	return nil
}
