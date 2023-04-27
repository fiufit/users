package repositories

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"go.uber.org/zap"
)

type Firebase interface {
	Register(ctx context.Context, email string, pw string) (string, error)
	DeleteUser(ctx context.Context, userID string) error
}

type FirebaseRepository struct {
	logger *zap.Logger
	auth   *auth.Client
}

func NewFirebaseRepository(logger *zap.Logger, auth *auth.Client) FirebaseRepository {
	return FirebaseRepository{logger: logger, auth: auth}
}

func (repo FirebaseRepository) DeleteUser(ctx context.Context, userID string) error {
	return repo.auth.DeleteUser(ctx, userID)
}

func (repo FirebaseRepository) Register(ctx context.Context, email string, pw string) (string, error) {
	user, err := repo.auth.GetUserByEmail(ctx, email)
	if err == nil && user != nil {
		if user.EmailVerified {
			return "", contracts.ErrUserAlreadyExists
		}

		updateUserParams := (&auth.UserToUpdate{}).Password(pw)
		updatedUser, err := repo.auth.UpdateUser(ctx, user.UID, updateUserParams)
		if err != nil {
			return "", err
		}
		return updatedUser.UID, nil
	}

	params := (&auth.UserToCreate{}).Email(email).Password(pw).EmailVerified(false)
	newUser, err := repo.auth.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}
	return newUser.UID, nil
}
