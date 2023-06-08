package accounts

import (
	"context"
	"time"

	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/contracts/metrics"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
)

type Registerer interface {
	Register(ctx context.Context, req accounts.RegisterRequest) (accounts.RegisterResponse, error)
	FinishRegister(ctx context.Context, req accounts.FinishRegisterRequest) (accounts.FinishRegisterResponse, error)
}

type RegistererImpl struct {
	users   repositories.Users
	metrics repositories.Metrics
	logger  *zap.Logger
	auth    repositories.Firebase
}

func NewRegisterImpl(users repositories.Users, logger *zap.Logger, auth repositories.Firebase, metrics repositories.Metrics) RegistererImpl {
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
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		Interests:         req.Interests,
	}
	createdUser, err := uc.users.CreateUser(ctx, usr)
	if err != nil {
		return accounts.FinishRegisterResponse{}, err
	}

	metricReq := metrics.CreateMetricRequest{
		MetricType: "register",
		SubType:    req.Method,
	}
	uc.metrics.Create(ctx, metricReq)

	return accounts.FinishRegisterResponse{User: createdUser}, nil
}
