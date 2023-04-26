package users

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserDoesNotExistError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	userID := "h0l4"

	req := users.UpdateUserRequest{
		ID:       userID,
		Nickname: "Arnold",
	}

	//when
	userRepo.On("GetByID", ctx, userID).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo)
	_, err := userUc.UpdateUser(ctx, req)
	assert.Error(t, err)
}

func TestUpdateUserPatchModelError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
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
	userRepo.On("GetByID", ctx, userID).Return(outdatedUser, nil)
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(otherUser, nil)
	userUc := NewUserUpdaterImpl(userRepo)
	_, err := userUc.UpdateUser(ctx, req)

	//then
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyExists)
}

func TestUpdateUserRepoError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
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
	userRepo.On("GetByID", ctx, userID).Return(outdatedUser, nil)
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo)
	patchedUser, _ := userUc.patchUserModel(ctx, outdatedUser, req)
	userRepo.On("Update", ctx, patchedUser).Return(models.User{}, errors.New("repo error"))
	_, err := userUc.UpdateUser(ctx, req)

	//then
	assert.Error(t, err)
}

func TestUpdateUserOk(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
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
	userRepo.On("GetByID", ctx, userID).Return(outdatedUser, nil)
	userRepo.On("GetByNickname", ctx, req.Nickname).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserUpdaterImpl(userRepo)
	patchedUser, _ := userUc.patchUserModel(ctx, outdatedUser, req)
	userRepo.On("Update", ctx, patchedUser).Return(models.User{}, nil)
	_, err := userUc.UpdateUser(ctx, req)

	//then
	assert.NoError(t, err)
}
