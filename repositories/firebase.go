package repositories

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"go.uber.org/zap"
)

type Firebase interface {
	Register(ctx context.Context, email string, pw string) (*auth.UserRecord, error)
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

func (repo FirebaseRepository) Register(ctx context.Context, email string, pw string) (*auth.UserRecord, error) {
	user, err := repo.auth.GetUserByEmail(ctx, email)
	if err == nil && user != nil {
		if user.EmailVerified {
			return nil, contracts.ErrUserAlreadyExists
		}

		updateUserParams := (&auth.UserToUpdate{}).Password(pw)
		updatedUser, err := repo.auth.UpdateUser(ctx, user.UID, updateUserParams)
		if err != nil {
			return nil, err
		}
		return updatedUser, nil
	}

	params := (&auth.UserToCreate{}).Email(email).Password(pw).EmailVerified(false)
	return repo.auth.CreateUser(ctx, params)
}
