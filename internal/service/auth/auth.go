package auth

import (
	"atlogex/gofoyer/internal/domain/models"
	"atlogex/gofoyer/internal/lib/jwt"
	"atlogex/gofoyer/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

type UserProvider interface {
	User(
		ctx context.Context, email string,
	) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (token string, err error) {
	const operation = "auth.LoginUser"

	log := a.log.With(slog.String("operation", operation), slog.String("email", email))
	log.Info("Login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found", err)
		}

		return "", fmt.Errorf("failed to get user %s: %w", operation, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("Invalid credentials", err)

		return "", fmt.Errorf("failed to compare password %s: %w", operation, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		a.log.Warn("App not found", err)

		return "", fmt.Errorf("failed to get app %s: %w", operation, err)
	}

	log.Info("User logged in", slog.Int64("appID", app.ID))

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("Failed to create token", err)

		return "", fmt.Errorf("failed to create token %s: %w", operation, err)
	}

	return token, nil
}

func (a Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const operation = "auth.RegisterNewUser"

	log := a.log.With(slog.String("operation", operation), slog.String("email", email))
	log.Info("Registering new user")

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", operation, err)

		return 0, fmt.Errorf("failed to hash password %s: %w", operation, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, bcryptHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("User already exists", err)

			return 0, fmt.Errorf("failed to save user %s: %w", operation, ErrUserExists)
		}

		log.Error("Failed to save user", operation, err)

		return 0, fmt.Errorf("failed to save user %s: %w", operation, err)
	}

	log.Info("Registered new user", slog.Int64("userID", id))

	return id, nil
}

func (a Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const operation = "auth.IsAdmin"

	log := a.log.With(
		slog.String("operation", operation),
		slog.Int64("user_id", userID),
	)

	log.Info("Checking if user is admin")

	IsAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("App not found", err)

			return false, fmt.Errorf("failed to check if user is admin %s: %w", operation, ErrInvalidAppId)
		}
		log.Error("Failed to check if user is admin", err)

		return false, fmt.Errorf("failed to check if user is admin %s: %w", operation, err)
	}

	log.Info("User is admin", slog.Bool("is_admin", IsAdmin))

	return IsAdmin, nil

}
