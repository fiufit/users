package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiufit/users/contracts/accounts"
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

	firebaseRepo.On("Register", ctx, req).Return(uid, nil)
	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo)
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

	firebaseRepo.On("Register", ctx, req).Return("", errors.New("repo error"))

	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo)
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
	req := accounts.FinishRegisterRequest{
		UserID:       "123456789",
		Nickname:     "Nick Test",
		DisplayName:  "Name Test",
		IsMale:       newTrue(),
		BirthDate:    birthDate,
		Height:       180,
		Weight:       80,
		MainLocation: "Testland",
		Interests:    []string{"Testing"},
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
		MainLocation:      req.MainLocation,
		Interests:         nil,
	}
	firebaseRepo := new(mocks.Firebase)
	userRepo := new(mocks.Users)

	_, _ = mpatch.PatchMethod(time.Now, func() time.Time {
		return creationDate
	})
	userRepo.On("CreateUser", ctx, usr).Return(usr, nil)
	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo)
	res, err := registerUc.FinishRegister(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.User, usr)
}

func TestFinishRegisterError(t *testing.T) {

	ctx := context.Background()
	birthDate := time.Now()
	req := accounts.FinishRegisterRequest{
		UserID:       "123456789",
		Nickname:     "Nick Test",
		DisplayName:  "Name Test",
		IsMale:       newTrue(),
		BirthDate:    birthDate,
		Height:       180,
		Weight:       80,
		MainLocation: "Testland",
		Interests:    []string{"Testing"},
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
		MainLocation:      req.MainLocation,
		Interests:         nil,
	}

	firebaseRepo := new(mocks.Firebase)
	userRepo := new(mocks.Users)

	_, _ = mpatch.PatchMethod(time.Now, func() time.Time {
		return creationDate
	})

	userRepo.On("CreateUser", ctx, usr).Return(models.User{}, errors.New("repo error"))
	registerUc := NewRegisterImpl(userRepo, zaptest.NewLogger(t), firebaseRepo)
	res, err := registerUc.FinishRegister(ctx, req)

	assert.Equal(t, res.User, models.User{})
	assert.Error(t, err)
}
