package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/contracts/metrics"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
	"go.uber.org/zap/zaptest"
)

func TestRegisterOk(t *testing.T) {

	ctx := context.Background()
	uid := "123456789"
	req := accounts.RegisterRequest{
		Email:    "test@fiufit.com",
		Password: "password",
	}
	firebaseRepo := new(mocks.Firebase)
	userRepo := new(mocks.Users)
	metricsRepo := new(mocks.Metrics)

	firebaseRepo.On("Register", ctx, req).Return(uid, nil)
	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo, metricsRepo)
	res, err := registerUc.Register(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.UserID, uid)
}

func TestRegisterError(t *testing.T) {

	ctx := context.Background()
	req := accounts.RegisterRequest{
		Email:    "test@fiufit.com",
		Password: "password",
	}
	firebaseRepo := new(mocks.Firebase)
	userRepo := new(mocks.Users)
	metricsRepo := new(mocks.Metrics)

	firebaseRepo.On("Register", ctx, req).Return("", errors.New("repo error"))

	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo, metricsRepo)
	res, err := registerUc.Register(ctx, req)

	assert.Equal(t, res.UserID, "")
	assert.Error(t, err)
}

func newTrue() *bool {
	b := true
	return &b
}

func TestFinishRegisterOk(t *testing.T) {

	ctx := context.Background()
	birthDate := time.Now()

	lat := 50.1
	long := 50.2

	req := accounts.FinishRegisterRequest{
		UserID:      "123456789",
		Nickname:    "Nick Test",
		DisplayName: "Name Test",
		IsMale:      newTrue(),
		BirthDate:   birthDate,
		Height:      180,
		Weight:      80,
		Latitude:    &lat,
		Longitude:   &long,
		Interests:   []models.Interest{},
	}
	creationDate := time.Now()
	usr := models.User{
		ID:                req.UserID,
		Nickname:          req.Nickname,
		DisplayName:       req.DisplayName,
		IsMale:            *req.IsMale,
		CreatedAt:         creationDate,
		BornAt:            req.BirthDate,
		Height:            req.Height,
		Weight:            req.Weight,
		IsVerifiedTrainer: false,
		Latitude:          lat,
		Longitude:         long,
		Interests:         req.Interests,
	}
	firebaseRepo := new(mocks.Firebase)
	userRepo := new(mocks.Users)
	metricsRepo := new(mocks.Metrics)

	registerMetricReq := metrics.CreateMetricRequest{
		MetricType: "register",
		SubType:    req.Method,
	}
	metricsRepo.On("Create", ctx, registerMetricReq)

	locationMetricsReq := metrics.CreateMetricRequest{
		MetricType: "location",
		SubType:    usr.MainLocation,
	}
	metricsRepo.On("Create", ctx, locationMetricsReq)

	timePatch, _ := mpatch.PatchMethod(time.Now, func() time.Time {
		return creationDate
	})
	userRepo.On("CreateUser", ctx, usr).Return(usr, nil)
	firebaseRepo.On("GetUserPictureUrl", ctx, usr.ID).Return("")
	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo, metricsRepo)
	_, err := registerUc.FinishRegister(ctx, req)

	timePatch.Unpatch()
	assert.NoError(t, err)
}

func TestFinishRegisterError(t *testing.T) {

	lat := 50.1
	long := 50.2
	ctx := context.Background()
	birthDate := time.Now()
	req := accounts.FinishRegisterRequest{
		UserID:      "123456789",
		Nickname:    "Nick Test",
		DisplayName: "Name Test",
		IsMale:      newTrue(),
		BirthDate:   birthDate,
		Height:      180,
		Weight:      80,
		Latitude:    &lat,
		Longitude:   &long,
		Interests:   []models.Interest{},
	}
	creationDate := time.Now()
	usr := models.User{
		ID:                req.UserID,
		Nickname:          req.Nickname,
		DisplayName:       req.DisplayName,
		IsMale:            *req.IsMale,
		CreatedAt:         creationDate,
		BornAt:            req.BirthDate,
		Height:            req.Height,
		Weight:            req.Weight,
		IsVerifiedTrainer: false,
		Latitude:          lat,
		Longitude:         long,
		Interests:         req.Interests,
	}

	firebaseRepo := new(mocks.Firebase)
	userRepo := new(mocks.Users)
	metricsRepo := new(mocks.Metrics)

	timePatch, _ := mpatch.PatchMethod(time.Now, func() time.Time {
		return creationDate
	})

	userRepo.On("CreateUser", ctx, usr).Return(models.User{}, errors.New("repo error"))
	firebaseRepo.On("GetUserPictureUrl", ctx, usr.ID).Return("")
	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo, metricsRepo)
	res, err := registerUc.FinishRegister(ctx, req)

	timePatch.Unpatch()
	assert.Equal(t, res.User, models.User{})
	assert.Error(t, err)
}
