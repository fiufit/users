// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/fiufit/users/models"
	mock "github.com/stretchr/testify/mock"
)

// VerificationPins is an autogenerated mock type for the VerificationPins type
type VerificationPins struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, pin
func (_m *VerificationPins) Create(ctx context.Context, pin models.VerificationPin) (models.VerificationPin, error) {
	ret := _m.Called(ctx, pin)

	var r0 models.VerificationPin
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.VerificationPin) (models.VerificationPin, error)); ok {
		return rf(ctx, pin)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.VerificationPin) models.VerificationPin); ok {
		r0 = rf(ctx, pin)
	} else {
		r0 = ret.Get(0).(models.VerificationPin)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.VerificationPin) error); ok {
		r1 = rf(ctx, pin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *VerificationPins) GetByUserID(ctx context.Context, userID string) (models.VerificationPin, error) {
	ret := _m.Called(ctx, userID)

	var r0 models.VerificationPin
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.VerificationPin, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.VerificationPin); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(models.VerificationPin)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewVerificationPins interface {
	mock.TestingT
	Cleanup(func())
}

// NewVerificationPins creates a new instance of VerificationPins. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVerificationPins(t mockConstructorTestingTNewVerificationPins) *VerificationPins {
	mock := &VerificationPins{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
