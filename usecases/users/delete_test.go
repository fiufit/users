package users

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserDeleter_DeleteUser_Error(t *testing.T) {
	usersMock := new(mocks.Users)
	uc := NewUserDeleterImpl(usersMock)
	ctx := context.Background()
	testUserID := "testUserID"
	usersMock.On("DeleteUser", ctx, testUserID).Return(errors.New("repo error"))

	err := uc.DeleteUser(ctx, testUserID)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "repo error")
}

func TestUserDeleter_DeleteUser_Ok(t *testing.T) {
	usersMock := new(mocks.Users)
	uc := NewUserDeleterImpl(usersMock)
	ctx := context.Background()
	testUserID := "testUserID"
	usersMock.On("DeleteUser", ctx, testUserID).Return(nil)

	err := uc.DeleteUser(ctx, testUserID)
	assert.NoError(t, err)
}
