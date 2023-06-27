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
	GetClosestUsers(ctx context.Context, req users.GetClosestUsersRequest) (users.GetUsersResponse, error)
	GetUserByNickname(ctx context.Context, nickname string) (models.User, error)
	GetUserFollowers(ctx context.Context, req users.GetUserFollowersRequest) (users.GetUserFollowersResponse, error)
	GetUserFollowed(ctx context.Context, req users.GetFollowedUsersRequest) (users.GetFollowedUsersResponse, error)
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
		return user, err
	}
	return user, nil
}

func (uc *UserGetterImpl) GetUserByNickname(ctx context.Context, nickname string) (models.User, error) {
	user, err := uc.users.GetByNickname(ctx, nickname)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (uc *UserGetterImpl) GetUsers(ctx context.Context, req users.GetUsersRequest) (users.GetUsersResponse, error) {
	res, err := uc.users.Get(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc *UserGetterImpl) GetClosestUsers(ctx context.Context, req users.GetClosestUsersRequest) (users.GetUsersResponse, error) {
	usr, err := uc.users.GetByID(ctx, req.UserID)
	if err != nil {
		return users.GetUsersResponse{}, err
	}
	req.Latitude = usr.Latitude
	req.Longitude = usr.Longitude

	res, err := uc.users.GetByDistance(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc *UserGetterImpl) GetUserFollowers(ctx context.Context, req users.GetUserFollowersRequest) (users.GetUserFollowersResponse, error) {
	res, err := uc.users.GetFollowers(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc *UserGetterImpl) GetUserFollowed(ctx context.Context, req users.GetFollowedUsersRequest) (users.GetFollowedUsersResponse, error) {
	res, err := uc.users.GetFollowed(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
