package tasks

import (
	"auth-sso/internal/tasks/handlers/identity"
	"github.com/hibiken/asynq"
)

func SetupTaskHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(identity.TaskIdentifier, identity.HandleIdentityVerificationTask)
}
