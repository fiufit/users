package users

import (
	"context"
	"testing"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestEnableUserOk(t *testing.T) {

	ctx := context.Background()
	uid := "123456789"
	user := models.User{ID: uid}
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	firebaseRepo.On("EnableUser", ctx, uid).Return(nil)
	userRepo.On("GetByID", ctx, uid).Return(user, nil)
	user.Disabled = false
	userRepo.On("Update", ctx, user).Return(user, nil)
	enableUserUc := NewUserEnablerImpl(userRepo, firebaseRepo, new(mocks.Metrics), zaptest.NewLogger(t))
	err := enableUserUc.EnableUser(ctx, uid)

	assert.NoError(t, err)
}

func TestEnableUserError(t *testing.T) {

	ctx := context.Background()
	uid := "123456789"
	user := models.User{ID: uid}
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	firebaseRepo.On("EnableUser", ctx, uid).Return(contracts.ErrUserNotDisabled)
	userRepo.On("GetByID", ctx, uid).Return(user, nil)
	enableUserUc := NewUserEnablerImpl(userRepo, firebaseRepo, new(mocks.Metrics), zaptest.NewLogger(t))
	err := enableUserUc.EnableUser(ctx, uid)

	assert.Error(t, err)
}

func TestEnableUserNotFoundError(t *testing.T) {

	ctx := context.Background()
	uid := "notFound"
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	userRepo.On("GetByID", ctx, uid).Return(models.User{}, contracts.ErrUserNotFound)
	enableUserUc := NewUserEnablerImpl(userRepo, firebaseRepo, new(mocks.Metrics), zaptest.NewLogger(t))
	err := enableUserUc.EnableUser(ctx, uid)

	assert.Error(t, err)
}

func TestDisableUserOk(t *testing.T) {

	ctx := context.Background()
	uid := "123456789"
	user := models.User{ID: uid}
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	firebaseRepo.On("DisableUser", ctx, uid).Return(nil)
	userRepo.On("GetByID", ctx, uid).Return(user, nil)
	user.Disabled = true
	userRepo.On("Update", ctx, user).Return(user, nil)
	enableUserUc := NewUserEnablerImpl(userRepo, firebaseRepo, new(mocks.Metrics), zaptest.NewLogger(t))
	err := enableUserUc.DisableUser(ctx, uid)

	assert.NoError(t, err)
}

func TestDisableUserError(t *testing.T) {

	ctx := context.Background()
	uid := "123456789"
	user := models.User{ID: uid}
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	firebaseRepo.On("DisableUser", ctx, uid).Return(contracts.ErrUserAlreadyDisabled)
	userRepo.On("GetByID", ctx, uid).Return(user, nil)
	enableUserUc := NewUserEnablerImpl(userRepo, firebaseRepo, new(mocks.Metrics), zaptest.NewLogger(t))
	err := enableUserUc.DisableUser(ctx, uid)

	assert.Error(t, err)
}

func TestDisableUserNotFoundError(t *testing.T) {

	ctx := context.Background()
	uid := "notFound"
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	userRepo.On("GetByID", ctx, uid).Return(models.User{}, contracts.ErrUserNotFound)
	enableUserUc := NewUserEnablerImpl(userRepo, firebaseRepo, new(mocks.Metrics), zaptest.NewLogger(t))
	err := enableUserUc.DisableUser(ctx, uid)

	assert.Error(t, err)
}
