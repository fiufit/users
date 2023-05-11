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
	users    repositories.Users
	firebase repositories.Firebase
	logger   *zap.Logger
}

func NewUserGetterImpl(users repositories.Users, firebase repositories.Firebase, logger *zap.Logger) UserGetterImpl {
	return UserGetterImpl{users: users, firebase: firebase, logger: logger}
}

func (uc *UserGetterImpl) GetUserByID(ctx context.Context, uid string) (models.User, error) {
	user, err := uc.users.GetByID(ctx, uid)
	if err != nil {
		return user, err
	}
	uc.fillUserPicture(ctx, &user)
	return user, nil
}

func (uc *UserGetterImpl) GetUserByNickname(ctx context.Context, nickname string) (models.User, error) {
	user, err := uc.users.GetByNickname(ctx, nickname)
	if err != nil {
		return user, err
	}
	uc.fillUserPicture(ctx, &user)
	return user, nil
}

func (uc *UserGetterImpl) GetUsers(ctx context.Context, req users.GetUsersRequest) (users.GetUsersResponse, error) {
	res, err := uc.users.Get(ctx, req)
	if err != nil {
		return res, err
	}
	for i := range res.Users {
		uc.fillUserPicture(ctx, &res.Users[i])
	}
	return res, nil
}

func (uc *UserGetterImpl) fillUserPicture(ctx context.Context, user *models.User) {
	userPictureUrl := uc.firebase.GetUserPictureUrl(ctx, user.ID)
	(*user).PictureUrl = userPictureUrl
}
