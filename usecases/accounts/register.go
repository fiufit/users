package accounts

import (
	"context"
	"errors"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/utils"
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
	mailer utils.Mailer
}

func NewRegisterImpl(users repositories.Users, logger *zap.Logger, auth *auth.Client, mailer utils.Mailer) RegistererImpl {
	return RegistererImpl{users: users, logger: logger, auth: auth, mailer: mailer}
}

func (uc *RegistererImpl) Register(ctx context.Context, req contracts.RegisterRequest) (contracts.RegisterResponse, error) {
	params := (&auth.UserToCreate{}).Email(req.Email).Password(req.Password).EmailVerified(false)
	user, err := uc.auth.CreateUser(ctx, params)
	if err != nil {
		return contracts.RegisterResponse{}, contracts.ErrUserAlreadyExists
	}

	verificationLink, err := uc.auth.EmailVerificationLink(ctx, user.Email)
	if err != nil {
		uc.logger.Error("Unable to generate verification link for email", zap.String("email", req.Email), zap.Error(err))
		return contracts.RegisterResponse{}, err
	}

	/*TODO Find out how to user firebase's email verification instead of our own mail account + SMTP server. Apparently
	the Go firebase SDK doesn't have auth.SendEmailVerification()
	*/
	err = uc.mailer.SendAccountVerificationEmail(user.Email, verificationLink)
	if err != nil {
		return contracts.RegisterResponse{}, err
	}
	return contracts.RegisterResponse{UserID: user.UID}, nil
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
