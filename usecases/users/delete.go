package users

import (
	"context"

	"github.com/fiufit/users/repositories"
)

type UserDeleter interface {
	DeleteUser(ctx context.Context, userID string) error
}

type UserDeleterImpl struct {
	users repositories.Users
}

func NewUserDeleterImpl(users repositories.Users) UserDeleterImpl {
	return UserDeleterImpl{users: users}
}

func (uc *UserDeleterImpl) DeleteUser(ctx context.Context, userID string) error {
	return uc.users.DeleteUser(ctx, userID)
}
