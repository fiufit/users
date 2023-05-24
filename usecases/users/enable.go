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
	firebase repositories.Firebase
	logger   *zap.Logger
}

func NewUserEnablerImpl(firebase repositories.Firebase, logger *zap.Logger) UserEnablerImpl {
	return UserEnablerImpl{firebase: firebase, logger: logger}
}

func (uc UserEnablerImpl) EnableUser(ctx context.Context, userID string) error {
	return uc.firebase.EnableUser(ctx, userID)
}

func (uc UserEnablerImpl) DisableUser(ctx context.Context, userID string) error {
	return uc.firebase.DisableUser(ctx, userID)
}
