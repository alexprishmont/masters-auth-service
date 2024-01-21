package auth

import (
	"auth-sso/internal/domain/models"
	"auth-sso/internal/storage"
	"auth-sso/lib/jwt"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log                *slog.Logger
	userSaver          UserSaver
	userProvider       UserProvider
	appProvider        AppProvider
	permissionProvider PermissionProvider
	tokenTTL           time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte,
	) (uid string, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type PermissionProvider interface {
	Can(ctx context.Context, permission string, userId string) (bool, error)
}

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorUserExists         = errors.New("user exists")
	ErrorAppNotFound        = errors.New("wrong application AppID")
	ErrorUserNotAuthorized  = errors.New("user action is not authorized")
)

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	permissionProvider PermissionProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:                log,
		userSaver:          userSaver,
		userProvider:       userProvider,
		appProvider:        appProvider,
		permissionProvider: permissionProvider,
		tokenTTL:           tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("Logging user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrorAppNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrorAppNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user AppID
// If user with given email address already exists, returns error.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("Registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("Failed to generate password hash", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passwordHash)

	if err != nil {
		if errors.Is(err, storage.ErrorUserExists) {
			log.Warn("user already exists", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrorUserExists)
		}

		log.Error("Failed to save user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User registered")

	return id, nil
}

func (a *Auth) Authorize(ctx context.Context, permission string, userId string) (isAuthorized bool, err error) {
	const op = "auth.Authorize"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("Authorizing user action")

	can, err := a.permissionProvider.Can(ctx, permission, userId)

	if err != nil {
		log.Error("failed to authorize user", slog.String("error", err.Error()))

		return false, fmt.Errorf("%s: %w", op, ErrorUserNotAuthorized)
	}

	return can, nil
}
