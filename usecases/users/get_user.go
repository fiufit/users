package users

import (
	"context"

	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
)

type UserGetter interface {
	GetUserByID(ctx context.Context, uid string) (models.User, error)
	GetUsers(ctx context.Context, req users.GetUsersRequest) (users.GetUsersResponse, error)
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
	return uc.users.GetByID(ctx, uid)
}

func (uc *UserGetterImpl) GetUserByNickname(ctx context.Context, nickname string) (models.User, error) {
	return uc.users.GetByNickname(ctx, nickname)
}

func (uc *UserGetterImpl) GetUsers(ctx context.Context, req users.GetUsersRequest) (users.GetUsersResponse, error) {
	return uc.users.Get(ctx, req)
}
