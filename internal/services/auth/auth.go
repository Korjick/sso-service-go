package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Korjick/sso-service-go/internal/domain/model"
	"github.com/Korjick/sso-service-go/internal/lib/jwt"
	"github.com/Korjick/sso-service-go/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	emptyUserIDValue = -1
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalError      = errors.New("internal error")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user model.User, err error)
	IsAdmin(ctx context.Context, uid int64) (isAdmin bool, err error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (app model.App, err error)
}

type Auth struct {
	logger       *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func New(
	logger *slog.Logger,
	saver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		logger:       logger,
		userSaver:    saver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int) (token string, err error) {
	const op = "auth.Login"

	logger := a.logger.With(slog.String("op", op))
	logger.Info("logging in user", slog.String("email", email))

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			logger.Error("user not found", slog.String("err", err.Error()))
			return "", fmt.Errorf("%s: user not found: %w", op, ErrInvalidCredentials)
		}
		logger.Error("failed to get user", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: failed to get user: %w", op, ErrInternalError)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.logger.Error("invalid credentials", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: invalid credentials: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrAppNotFound):
			logger.Error("app not found", slog.String("err", err.Error()))
			return "", fmt.Errorf("%s: app not found: %w", op, ErrInvalidAppID)
		}
		logger.Error("failed to get app", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: failed to get app: %w", op, ErrInternalError)
	}

	logger.Info("user is logged in", slog.String("email", email))

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		logger.Error("failed to generate token", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: failed to generate token: %w", op, ErrInternalError)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	logger := a.logger.With(slog.String("op", op))
	logger.Info("registering new user", slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", slog.String("err", err.Error()))
		return emptyUserIDValue, fmt.Errorf("%s: failed to generate password hash: %w", op, ErrInternalError)
	}

	uid, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserExists):
			logger.Error("user exists", slog.String("err", err.Error()))
			return emptyUserIDValue, fmt.Errorf("%s: user exists: %w", op, ErrUserAlreadyExists)
		}
		logger.Error("failed to save user", slog.String("err", err.Error()))
		return emptyUserIDValue, fmt.Errorf("%s: failed to save user: %w", op, ErrInternalError)
	}

	logger.Info("user is registered", slog.Int64("uid", uid))
	return uid, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	const op = "auth.IsAdmin"

	logger := a.logger.With(slog.String("op", op))
	logger.Info("checking if user is admin", slog.Int64("uid", userID))

	isAdmin, err = a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			logger.Error("user not found", slog.String("err", err.Error()))
			return false, fmt.Errorf("%s: user not found: %w", op, ErrInvalidCredentials)
		}
		logger.Error("failed to get user", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: failed to get user: %w", op, ErrInternalError)
	}

	return isAdmin, nil
}
