package accounts

import (
	"context"
	"time"

	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/contracts/metrics"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/repositories/external"
	"go.uber.org/zap"
)

type Registerer interface {
	Register(ctx context.Context, req accounts.RegisterRequest) (accounts.RegisterResponse, error)
	FinishRegister(ctx context.Context, req accounts.FinishRegisterRequest) (accounts.FinishRegisterResponse, error)
}

type RegistererImpl struct {
	users   repositories.Users
	metrics external.Metrics
	logger  *zap.Logger
	auth    external.Firebase
}

func NewRegisterImpl(users repositories.Users, logger *zap.Logger, auth external.Firebase, metrics external.Metrics) RegistererImpl {
	return RegistererImpl{users: users, metrics: metrics, logger: logger, auth: auth}
}

func (uc *RegistererImpl) Register(ctx context.Context, req accounts.RegisterRequest) (accounts.RegisterResponse, error) {
	newUserID, err := uc.auth.Register(ctx, req)
	if err != nil {
		return accounts.RegisterResponse{}, err
	}

	return accounts.RegisterResponse{UserID: newUserID}, nil
}

func (uc *RegistererImpl) FinishRegister(ctx context.Context, req accounts.FinishRegisterRequest) (accounts.FinishRegisterResponse, error) {

	usr := models.User{
		ID:                req.UserID,
		Nickname:          req.Nickname,
		DisplayName:       req.DisplayName,
		IsMale:            *req.IsMale,
		CreatedAt:         time.Now(),
		BornAt:            req.BirthDate,
		Height:            req.Height,
		Weight:            req.Weight,
		IsVerifiedTrainer: false,
		Latitude:          *req.Latitude,
		Longitude:         *req.Longitude,
		Interests:         req.Interests,
	}
	createdUser, err := uc.users.CreateUser(ctx, usr)
	if err != nil {
		return accounts.FinishRegisterResponse{}, err
	}

	createdUser.PictureUrl = uc.auth.GetUserPictureUrl(ctx, createdUser.ID)

	registerMetricReq := metrics.CreateMetricRequest{
		MetricType: "register",
		SubType:    req.Method,
	}
	uc.metrics.Create(ctx, registerMetricReq)

	locationMetricReq := metrics.CreateMetricRequest{
		MetricType: "location",
		SubType:    usr.MainLocation,
	}
	uc.metrics.Create(ctx, locationMetricReq)

	return accounts.FinishRegisterResponse{User: createdUser}, nil
}
