// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/fiufit/users/models"
	mock "github.com/stretchr/testify/mock"

	users "github.com/fiufit/users/contracts/users"
)

// Users is an autogenerated mock type for the Users type
type Users struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *Users) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *Users) DeleteUser(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FollowUser provides a mock function with given fields: ctx, followedUserID, followerUserID
func (_m *Users) FollowUser(ctx context.Context, followedUserID string, followerUserID string) error {
	ret := _m.Called(ctx, followedUserID, followerUserID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, followedUserID, followerUserID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, req
func (_m *Users) Get(ctx context.Context, req users.GetUsersRequest) (users.GetUsersResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 users.GetUsersResponse
	if rf, ok := ret.Get(0).(func(context.Context, users.GetUsersRequest) users.GetUsersResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(users.GetUsersResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, users.GetUsersRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, userID
func (_m *Users) GetByID(ctx context.Context, userID string) (models.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNickname provides a mock function with given fields: ctx, nickname
func (_m *Users) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	ret := _m.Called(ctx, nickname)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(ctx, nickname)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, nickname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFollowers provides a mock function with given fields: ctx, request
func (_m *Users) GetFollowers(ctx context.Context, request users.GetUserFollowersRequest) (users.GetUserFollowersResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 users.GetUserFollowersResponse
	if rf, ok := ret.Get(0).(func(context.Context, users.GetUserFollowersRequest) users.GetUserFollowersResponse); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(users.GetUserFollowersResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, users.GetUserFollowersRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnfollowUser provides a mock function with given fields: ctx, followedUserID, followerUserID
func (_m *Users) UnfollowUser(ctx context.Context, followedUserID string, followerUserID string) error {
	ret := _m.Called(ctx, followedUserID, followerUserID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, followedUserID, followerUserID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, user
func (_m *Users) Update(ctx context.Context, user models.User) (models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUsers interface {
	mock.TestingT
	Cleanup(func())
}

// NewUsers creates a new instance of Users. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUsers(t mockConstructorTestingTNewUsers) *Users {
	mock := &Users{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
