package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

//go:generate mockery --name Firebase
type Firebase interface {
	Register(ctx context.Context, req accounts.RegisterRequest) (string, error)
	DeleteUser(ctx context.Context, userID string) error
	GetUserPictureUrl(ctx context.Context, userID string) string
	EnableUser(ctx context.Context, userID string) error
	DisableUser(ctx context.Context, userID string) error
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

func (repo FirebaseRepository) DisableUser(ctx context.Context, userID string) error {
	usr, err := repo.auth.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if usr.Disabled {
		return contracts.ErrUserAlreadyDisabled
	}
	updateUserParams := (&auth.UserToUpdate{}).Disabled(true)
	_, err = repo.auth.UpdateUser(ctx, userID, updateUserParams)
	return err
}

func (repo FirebaseRepository) EnableUser(ctx context.Context, userID string) error {
	usr, err := repo.auth.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if !usr.Disabled {
		return contracts.ErrUserNotDisabled
	}
	updateUserParams := (&auth.UserToUpdate{}).Disabled(false)
	_, err = repo.auth.UpdateUser(ctx, userID, updateUserParams)
	return err
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
	defaultPicturePath := "profile_pictures/default.png"
	userPicturePath := "profile_pictures/" + userID + "/profile.png"
	picturePath := userPicturePath

	pictureHandle := repo.storageBucket.Object(userPicturePath)
	_, err := pictureHandle.Attrs(ctx)
	if err != nil {
		if !errors.Is(err, storage.ErrObjectNotExist) {
			repo.logger.Error("Unable to retrieve User picture from firebase storage", zap.String("userID", userID))
		}
		picturePath = defaultPicturePath
	}

	opts := storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Hour * 24),
	}
	pictureUrl, err := repo.storageBucket.SignedURL(picturePath, &opts)
	if err != nil {
		pictureUrl = ""
		repo.logger.Error("Unable to Sign user picture from firebase storage", zap.String("userID", userID))
	}
	return pictureUrl
}
