package users

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts"
	metrics2 "github.com/fiufit/users/contracts/metrics"
	uContracts "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestUserFollowerImpl_FollowUser_FollowedDoesNotExist(t *testing.T) {
	users := new(mocks.Users)
	metrics := new(mocks.Metrics)
	notifications := new(mocks.Notifications)
	logger := zaptest.NewLogger(t)
	userFollower := NewUserFollowerImpl(users, notifications, metrics, logger)
	ctx := context.Background()
	req := uContracts.FollowUserRequest{FollowerUserID: "a", FollowedUserID: "b"}
	users.On("GetByID", ctx, req.FollowedUserID).Return(models.User{}, contracts.ErrUserNotFound)

	err := userFollower.FollowUser(ctx, req)

	assert.Error(t, err)
}

func TestUserFollowerImpl_FollowUser_FollowerDoesNotExist(t *testing.T) {
	users := new(mocks.Users)
	metrics := new(mocks.Metrics)
	notifications := new(mocks.Notifications)
	logger := zaptest.NewLogger(t)
	userFollower := NewUserFollowerImpl(users, notifications, metrics, logger)
	ctx := context.Background()
	req := uContracts.FollowUserRequest{FollowerUserID: "a", FollowedUserID: "b"}
	users.On("GetByID", ctx, req.FollowedUserID).Return(models.User{}, nil)
	users.On("GetByID", ctx, req.FollowerUserID).Return(models.User{}, contracts.ErrUserNotFound)

	err := userFollower.FollowUser(ctx, req)

	assert.Error(t, err)
}

func TestUserFollowerImpl_FollowUser_RepoError(t *testing.T) {
	users := new(mocks.Users)
	metrics := new(mocks.Metrics)
	notifications := new(mocks.Notifications)
	logger := zaptest.NewLogger(t)
	userFollower := NewUserFollowerImpl(users, notifications, metrics, logger)
	ctx := context.Background()
	req := uContracts.FollowUserRequest{FollowerUserID: "a", FollowedUserID: "b"}
	users.On("GetByID", ctx, req.FollowedUserID).Return(models.User{}, nil)
	users.On("GetByID", ctx, req.FollowerUserID).Return(models.User{}, nil)
	users.On("FollowUser", ctx, models.User{}, models.User{}).Return(errors.New("repo error"))
	err := userFollower.FollowUser(ctx, req)

	assert.Error(t, err)
}

func TestUserFollowerImpl_FollowUser_Ok(t *testing.T) {
	users := new(mocks.Users)
	metrics := new(mocks.Metrics)
	notifications := new(mocks.Notifications)
	logger := zaptest.NewLogger(t)
	userFollower := NewUserFollowerImpl(users, notifications, metrics, logger)
	ctx := context.Background()
	req := uContracts.FollowUserRequest{FollowerUserID: "a", FollowedUserID: "b"}
	users.On("GetByID", ctx, req.FollowedUserID).Return(models.User{}, nil)
	users.On("GetByID", ctx, req.FollowerUserID).Return(models.User{}, nil)
	users.On("FollowUser", ctx, models.User{}, models.User{}).Return(nil)
	notifications.On("SendFollowersNotification", ctx, models.User{}, models.User{}).Return(nil)
	metrics.On("Create", ctx, metrics2.CreateMetricRequest{MetricType: "user_followed", SubType: ""}).Return(nil)
	err := userFollower.FollowUser(ctx, req)

	assert.NoError(t, err)
}
