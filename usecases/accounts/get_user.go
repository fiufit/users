package accounts

import (
	"context"
	"errors"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
)

type UserGetter interface {
	GetUserByID(ctx context.Context, uid string) (models.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (models.User, error)
}

type UserGetterImpl struct {
	users  repositories.Users
	logger *zap.Logger
}

func NewUserGetterImpl(users repositories.Users, logger *zap.Logger) UserGetterImpl {
	return UserGetterImpl{users: users, logger: logger}
}

func (uc *UserGetterImpl) GetUserByID(ctx context.Context, uid string) (models.User, error) {
	user, err := uc.users.GetByID(ctx, uid)
	if err != nil {
		if errors.Is(err, contracts.ErrUserNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (uc *UserGetterImpl) GetUserByNickname(ctx context.Context, nickname string) (models.User, error) {
	user, err := uc.users.GetByNickname(ctx, nickname)
	if err != nil {
		if errors.Is(err, contracts.ErrUserNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}
