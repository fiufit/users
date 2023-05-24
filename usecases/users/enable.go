package users

import (
	"context"

	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
)

type UserEnabler interface {
	EnableUser(ctx context.Context, userID string) error
	DisableUser(ctx context.Context, userID string) error
}

type UserEnablerImpl struct {
	users    repositories.Users
	firebase repositories.Firebase
	logger   *zap.Logger
}

func NewUserEnablerImpl(users repositories.Users, firebase repositories.Firebase, logger *zap.Logger) UserEnablerImpl {
	return UserEnablerImpl{users: users, firebase: firebase, logger: logger}
}

func (uc UserEnablerImpl) EnableUser(ctx context.Context, userID string) error {
	_, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	return uc.firebase.EnableUser(ctx, userID)
}

func (uc UserEnablerImpl) DisableUser(ctx context.Context, userID string) error {
	_, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	return uc.firebase.DisableUser(ctx, userID)
}
