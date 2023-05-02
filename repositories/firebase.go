package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type Firebase interface {
	Register(ctx context.Context, req accounts.RegisterRequest) (string, error)
	DeleteUser(ctx context.Context, userID string) error
	GetUserPictureUrl(ctx context.Context, userID string) string
}

type FirebaseRepository struct {
	logger            *zap.Logger
	app               *firebase.App
	auth              *auth.Client
	storageBucketName string
	storageBucket     *storage.BucketHandle
}

func NewFirebaseRepository(logger *zap.Logger, sdkJson []byte, storageBucketName string) (FirebaseRepository, error) {
	opt := option.WithCredentialsJSON(sdkJson)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return FirebaseRepository{}, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		return FirebaseRepository{}, err
	}

	storageClient, err := app.Storage(context.Background())
	if err != nil {
		return FirebaseRepository{}, err
	}

	storageBucket, err := storageClient.Bucket(storageBucketName)
	if err != nil {
		return FirebaseRepository{}, err
	}

	repo := FirebaseRepository{
		logger:            logger,
		app:               app,
		auth:              auth,
		storageBucketName: storageBucketName,
		storageBucket:     storageBucket,
	}

	return repo, nil
}

func (repo FirebaseRepository) DeleteUser(ctx context.Context, userID string) error {
	return repo.auth.DeleteUser(ctx, userID)
}

func (repo FirebaseRepository) Register(ctx context.Context, req accounts.RegisterRequest) (string, error) {
	email := strings.ToLower(req.Email)
	pw := req.Password
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

func (repo FirebaseRepository) GetUserPictureUrl(ctx context.Context, userID string) string {
	pictureHandle := repo.storageBucket.Object("profile_pictures/" + userID + "/profile.png")
	pictureData, err := pictureHandle.Attrs(ctx)
	if err != nil {
		if !errors.Is(err, storage.ErrObjectNotExist) {
			repo.logger.Error("Unable to retrieve User picture from firebase storage", zap.String("userID", userID))
		}
		return fmt.Sprintf("https://storage.cloud.google.com/%v/profile_pictures/default.png", repo.storageBucketName)

	}
	return pictureData.MediaLink
}
