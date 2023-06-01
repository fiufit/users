package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserDoesNotExistError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)

	ctx := context.Background()
	userID := "h0l4"

	req := users.UpdateUserRequest{
		ID:       userID,
		Nickname: "Arnold",
	}

	//when
	firebaseRepo.On("GetUserPictureUrl", ctx, userID).Return("aHR0cHM6Ly9zaG9ydHVybC5hdC9mcHRXNg==")
	userRepo.On("GetByID", ctx, userID).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo, firebaseRepo)
	_, err := userUc.UpdateUser(ctx, req)
	assert.Error(t, err)
}

func TestUpdateUserPatchModelError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)
	ctx := context.Background()
	userID := "h0l4"

	req := users.UpdateUserRequest{
		ID:       userID,
		Nickname: "Arnold2",
	}

	outdatedUser := models.User{
		ID:       userID,
		Nickname: "Arnold",
	}

	otherUser := models.User{
		ID:       "blablala",
		Nickname: "Arnold2",
	}

	//when
	firebaseRepo.On("GetUserPictureUrl", ctx, userID).Return("aHR0cHM6Ly9zaG9ydHVybC5hdC9mcHRXNg==")
	userRepo.On("GetByID", ctx, userID).Return(outdatedUser, nil)
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(otherUser, nil)
	userUc := NewUserUpdaterImpl(userRepo, firebaseRepo)
	_, err := userUc.UpdateUser(ctx, req)

	//then
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyExists)
}

func TestUpdateUserRepoError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)
	ctx := context.Background()
	userID := "h0l4"

	req := users.UpdateUserRequest{
		ID:       userID,
		Nickname: "Arnold2",
	}

	outdatedUser := models.User{
		ID:       userID,
		Nickname: "Arnold",
	}

	//when
	firebaseRepo.On("GetUserPictureUrl", ctx, userID).Return("aHR0cHM6Ly9zaG9ydHVybC5hdC9mcHRXNg==")
	userRepo.On("GetByID", ctx, userID).Return(outdatedUser, nil)
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo, firebaseRepo)
	patchedUser, _ := userUc.patchUserModel(ctx, outdatedUser, req)
	userRepo.On("Update", ctx, patchedUser).Return(models.User{}, errors.New("repo error"))
	_, err := userUc.UpdateUser(ctx, req)

	//then
	assert.Error(t, err)
}

func TestUpdateUserOk(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)
	ctx := context.Background()
	userID := "h0l4"

	req := users.UpdateUserRequest{
		ID:       userID,
		Nickname: "Arnold2",
	}

	outdatedUser := models.User{
		ID:       userID,
		Nickname: "Arnold",
	}

	//when
	firebaseRepo.On("GetUserPictureUrl", ctx, userID).Return("aHR0cHM6Ly9zaG9ydHVybC5hdC9mcHRXNg==")
	userRepo.On("GetByID", ctx, userID).Return(outdatedUser, nil)
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo, firebaseRepo)
	patchedUser, _ := userUc.patchUserModel(ctx, outdatedUser, req)
	userRepo.On("Update", ctx, patchedUser).Return(models.User{}, nil)
	_, err := userUc.UpdateUser(ctx, req)

	//then
	assert.NoError(t, err)
}

func TestPatchUserModelOk(t *testing.T) {
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)
	ctx := context.Background()
	userID := "h0l4"

	outdatedUser := models.User{
		ID:                userID,
		Nickname:          "Arnold",
		DisplayName:       "Arnold Schwarzenegger",
		IsMale:            true,
		BornAt:            time.Now(),
		Height:            100,
		Weight:            100,
		IsVerifiedTrainer: false,
		Interests:         nil,
	}

	isMale := false
	req := users.UpdateUserRequest{
		ID:          userID,
		Nickname:    "Arnold2",
		DisplayName: "Arnie",
		IsMale:      &isMale,
		BirthDate:   time.Now().Add(1),
		Weight:      200,
		Height:      200,
	}

	firebaseRepo.On("GetUserPictureUrl", ctx, userID).Return("aHR0cHM6Ly9zaG9ydHVybC5hdC9mcHRXNg==")
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo, firebaseRepo)
	updatedUser, err := userUc.patchUserModel(ctx, outdatedUser, req)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Nickname, req.Nickname)
	assert.Equal(t, updatedUser.DisplayName, req.DisplayName)
	assert.Equal(t, updatedUser.IsMale, *req.IsMale)
	assert.Equal(t, updatedUser.BornAt, req.BirthDate)
	assert.Equal(t, updatedUser.Weight, req.Weight)
	assert.Equal(t, updatedUser.Height, req.Height)
}

func TestPatchUserModelRepoError(t *testing.T) {
	userRepo := new(mocks.Users)
	firebaseRepo := new(mocks.Firebase)
	ctx := context.Background()
	userID := "h0l4"

	outdatedUser := models.User{
		ID:       userID,
		Nickname: "Arnold",
	}

	req := users.UpdateUserRequest{
		ID:       userID,
		Nickname: "Arnold2",
	}

	firebaseRepo.On("GetUserPictureUrl", ctx, userID).Return("aHR0cHM6Ly9zaG9ydHVybC5hdC9mcHRXNg==")
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(models.User{}, errors.New("repo error"))
	userUc := NewUserUpdaterImpl(userRepo, firebaseRepo)
	_, err := userUc.patchUserModel(ctx, outdatedUser, req)

	assert.Error(t, err)
}
