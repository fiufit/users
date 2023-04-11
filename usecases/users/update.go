package users

import (
	"context"
	"errors"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
)

type UserUpdater interface {
	UpdateUser(ctx context.Context, req ucontracts.UpdateUserRequest) (models.User, error)
}

type UserUpdaterImpl struct {
	users repositories.Users
}

func NewUserUpdaterImpl(users repositories.Users) UserUpdaterImpl {
	return UserUpdaterImpl{users: users}
}

func (uc *UserUpdaterImpl) UpdateUser(ctx context.Context, req ucontracts.UpdateUserRequest) (models.User, error) {
	user, err := uc.users.GetByID(ctx, req.ID)
	if err != nil {
		return models.User{}, err
	}

	patchedUser, err := uc.patchUserModel(ctx, user, req)
	if err != nil {
		return models.User{}, err
	}

	updatedUser, err := uc.users.Update(ctx, patchedUser)
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}

func (uc *UserUpdaterImpl) patchUserModel(ctx context.Context, user models.User, req ucontracts.UpdateUserRequest) (models.User, error) {
	if req.Nickname != "" {
		_, err := uc.users.GetByNickname(ctx, req.Nickname)
		if err != nil && !errors.Is(err, contracts.ErrUserNotFound) {
			return models.User{}, err
		}
		if err == nil { // there is already a user with the desired nickname
			return models.User{}, contracts.ErrUserAlreadyExists
		}

		user.Nickname = req.Nickname
	}

	if req.IsMale != user.IsMale {
		user.IsMale = req.IsMale
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	if !req.BirthDate.IsZero() {
		user.BornAt = req.BirthDate
	}

	if req.Weight != 0 {
		user.Weight = req.Weight
	}

	if req.Height != 0 {
		user.Height = req.Height
	}

	if req.MainLocation != "" {
		user.MainLocation = req.MainLocation
	}

	return user, nil
}
