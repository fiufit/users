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
	usr, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	err = uc.firebase.EnableUser(ctx, userID)
	if err != nil {
		return err
	}
	usr.Disabled = false
	_, err = uc.users.Update(ctx, usr)
	if err != nil {
		uc.logger.Error("Unable to fully enable user", zap.Error(err), zap.Any("user", userID))
	}
	return err
}

func (uc UserEnablerImpl) DisableUser(ctx context.Context, userID string) error {
	usr, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	err = uc.firebase.DisableUser(ctx, userID)
	if err != nil {
		return err
	}
	usr.Disabled = true
	_, err = uc.users.Update(ctx, usr)
	if err != nil {
		uc.logger.Error("Unable to fully disable user", zap.Error(err), zap.Any("user", userID))
	}
	return err
}
