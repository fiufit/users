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
	"go.uber.org/zap/zaptest"
)

func TestGetUserByIdError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	userID := "H014"

	//when
	userRepo.On("GetByID", ctx, userID).Return(models.User{}, errors.New("repo error"))
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))
	_, err := userUc.GetUserByID(ctx, userID)

	//then
	assert.Error(t, err)
}

func TestGetUserByIdOk(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	userID := "H014"
	user := models.User{ID: userID}

	//when
	userRepo.On("GetByID", ctx, userID).Return(user, nil)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))
	res, err := userUc.GetUserByID(ctx, userID)

	//then
	assert.NoError(t, err)
	assert.Equal(t, res.ID, userID)
}

func TestGetUserByNicknameError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	username := "Arnold84"

	//when
	userRepo.On("GetByNickname", ctx, username).Return(models.User{}, errors.New("repo error"))
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))
	_, err := userUc.GetUserByNickname(ctx, username)

	//then
	assert.Error(t, err)
}

func TestGetUserByNicknameOk(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	username := "Arnold84"
	user := models.User{Nickname: username, ID: "userID"}

	//when
	userRepo.On("GetByNickname", ctx, username).Return(user, nil)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))
	res, err := userUc.GetUserByNickname(ctx, username)

	//then
	assert.NoError(t, err)
	assert.Equal(t, res.Nickname, username)
}

func TestGetUsersError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetUsersRequest{
		Name: "Arnold",
	}

	//when
	userRepo.On("Get", ctx, req).Return(users.GetUsersResponse{}, errors.New("repo error"))
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))
	_, err := userUc.GetUsers(ctx, req)

	//then
	assert.Error(t, err)
}

func TestGetUsersOk(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetUsersRequest{
		Name: "Arnold",
	}
	res := users.GetUsersResponse{
		Pagination: contracts.Pagination{},
		Users:      []models.User{},
	}

	//when
	userRepo.On("Get", ctx, req).Return(res, nil)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))
	_, err := userUc.GetUsers(ctx, req)

	//then
	assert.NoError(t, err)
}

func TestGetClosestUsers_ErrUserNotFound(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetClosestUsersRequest{
		UserID: "H014",
	}
	userRepo.On("GetByID", ctx, req.UserID).Return(models.User{}, contracts.ErrUserNotFound)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetClosestUsers(ctx, req)

	//then
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestGetClosestUsers_GetByDistanceError(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetClosestUsersRequest{
		UserID: "H014",
	}
	userRepo.On("GetByID", ctx, req.UserID).Return(models.User{ID: req.UserID}, nil)
	userRepo.On("GetByDistance", ctx, req).Return(users.GetUsersResponse{}, errors.New("repo error"))
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetClosestUsers(ctx, req)

	//then
	assert.Error(t, err)
}

func TestGetClosestUsers_Ok(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetClosestUsersRequest{
		UserID: "H014",
	}
	userRepo.On("GetByID", ctx, req.UserID).Return(models.User{ID: req.UserID}, nil)
	userRepo.On("GetByDistance", ctx, req).Return(users.GetUsersResponse{}, nil)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetClosestUsers(ctx, req)

	//then
	assert.NoError(t, err)
}

func TestGetUserFollowers_Error(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetUserFollowersRequest{
		UserID: "H014",
	}
	userRepo.On("GetFollowers", ctx, req).Return(users.GetUserFollowersResponse{}, errors.New("repo error"))
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetUserFollowers(ctx, req)

	//then
	assert.Error(t, err)
}

func TestGetUserFollowers_Ok(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetUserFollowersRequest{
		UserID: "H014",
	}
	userRepo.On("GetFollowers", ctx, req).Return(users.GetUserFollowersResponse{}, nil)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetUserFollowers(ctx, req)

	//then
	assert.NoError(t, err)
}

func TestGetUserFollowed_Error(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetFollowedUsersRequest{
		UserID: "H014",
	}
	userRepo.On("GetFollowed", ctx, req).Return(users.GetFollowedUsersResponse{}, errors.New("repo error"))
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetUserFollowed(ctx, req)

	//then
	assert.Error(t, err)
}

func TestGetUserFollowed_Ok(t *testing.T) {
	//given
	userRepo := new(mocks.Users)
	ctx := context.Background()
	req := users.GetFollowedUsersRequest{
		UserID: "H014",
	}
	userRepo.On("GetFollowed", ctx, req).Return(users.GetFollowedUsersResponse{}, nil)
	userUc := NewUserGetterImpl(userRepo, zaptest.NewLogger(t))

	//when
	_, err := userUc.GetUserFollowed(ctx, req)

	//then
	assert.NoError(t, err)
}
