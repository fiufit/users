package repositories

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"go.uber.org/zap"
)

type Firebase interface {
	Register(ctx context.Context, email string, pw string) (*auth.UserRecord, error)
	CreateUser(ctx context.Context, params *auth.UserToCreate) (*auth.UserRecord, error)
	UpdateUser(ctx context.Context, uid string, updateUserParams *auth.UserToUpdate) (*auth.UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
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

func (repo FirebaseRepository) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	return repo.auth.GetUserByEmail(ctx, email)
}

func (repo FirebaseRepository) UpdateUser(ctx context.Context, uid string, updateUserParams *auth.UserToUpdate) (*auth.UserRecord, error) {
	return repo.auth.UpdateUser(ctx, uid, updateUserParams)
}

func (repo FirebaseRepository) CreateUser(ctx context.Context, params *auth.UserToCreate) (*auth.UserRecord, error) {
	return repo.auth.CreateUser(ctx, params)
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
