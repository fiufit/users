package accounts

import (
	"context"
	"errors"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
)

type Registerer interface {
	Register(ctx context.Context, req contracts.RegisterRequest) (contracts.RegisterResponse, error)
	FinishRegister(ctx context.Context, req contracts.FinishRegisterRequest) (contracts.FinishRegisterResponse, error)
}

type RegistererImpl struct {
	users  repositories.Users
	logger *zap.Logger
	auth   *auth.Client
}

func NewRegisterImpl(users repositories.Users, logger *zap.Logger, auth *auth.Client) RegistererImpl {
	return RegistererImpl{users: users, logger: logger, auth: auth}
}

func (uc *RegistererImpl) Register(ctx context.Context, req contracts.RegisterRequest) (contracts.RegisterResponse, error) {

	user, err := uc.auth.GetUserByEmail(ctx, req.Email)
	if err == nil && user != nil {
		if user.EmailVerified {
			return contracts.RegisterResponse{}, contracts.ErrUserAlreadyExists
		}

		updateUserParams := (&auth.UserToUpdate{}).Password(req.Password)
		updatedUser, err := uc.auth.UpdateUser(ctx, user.UID, updateUserParams)
		if err != nil {
			return contracts.RegisterResponse{}, err
		}
		return contracts.RegisterResponse{UserID: updatedUser.UID}, nil
	}

	params := (&auth.UserToCreate{}).Email(req.Email).Password(req.Password).EmailVerified(false)
	newUser, err := uc.auth.CreateUser(ctx, params)
	if err != nil {
		return contracts.RegisterResponse{}, err
	}

	return contracts.RegisterResponse{UserID: newUser.UID}, nil
}

func (uc *RegistererImpl) FinishRegister(ctx context.Context, req contracts.FinishRegisterRequest) (contracts.FinishRegisterResponse, error) {
	_, err := uc.users.GetByID(ctx, req.UserID)
	if !errors.Is(err, contracts.ErrUserNotFound) {
		return contracts.FinishRegisterResponse{}, contracts.ErrUserAlreadyExists
	}

	usr := models.User{
		ID:                req.UserID,
		Nickname:          req.Nickname,
		DisplayName:       req.DisplayName,
		IsMale:            req.IsMale,
		CreatedAt:         time.Now(),
		BornAt:            req.BirthDate,
		Height:            req.Height,
		Weight:            req.Weight,
		IsVerifiedTrainer: false,
		MainLocation:      req.MainLocation,
		Interests:         nil,
	}
	_, err = uc.users.CreateUser(ctx, usr)
	return contracts.FinishRegisterResponse{User: usr}, err
}
