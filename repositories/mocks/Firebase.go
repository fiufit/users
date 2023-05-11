// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	accounts "github.com/fiufit/users/contracts/accounts"

	mock "github.com/stretchr/testify/mock"
)

// Firebase is an autogenerated mock type for the Firebase type
type Firebase struct {
	mock.Mock
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *Firebase) DeleteUser(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserPictureUrl provides a mock function with given fields: ctx, userID
func (_m *Firebase) GetUserPictureUrl(ctx context.Context, userID string) string {
	ret := _m.Called(ctx, userID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Register provides a mock function with given fields: ctx, req
func (_m *Firebase) Register(ctx context.Context, req accounts.RegisterRequest) (string, error) {
	ret := _m.Called(ctx, req)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, accounts.RegisterRequest) string); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, accounts.RegisterRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFirebase interface {
	mock.TestingT
	Cleanup(func())
}

// NewFirebase creates a new instance of Firebase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFirebase(t mockConstructorTestingTNewFirebase) *Firebase {
	mock := &Firebase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
