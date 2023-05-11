package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestUserRepository_CreateUser_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{}

	_ = repo.db.AddError(errors.New("test error"))

	_, err := repo.CreateUser(ctx, testUser)
	assert.Error(t, err)
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestUserRepository_CreateUser_DuplicatedError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test"}

	db.Create(&testUser)

	_, err := repo.CreateUser(ctx, testUser)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyExists)
}

func TestUserRepository_CreateUser_Ok(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test"}

	_, err := repo.CreateUser(ctx, testUser)

	var createdUser models.User
	db.First(&createdUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, createdUser.ID)
}

func TestUserRepository_GetByID_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test"}
	_ = db.AddError(errors.New("test error"))

	_, err := repo.GetByID(ctx, testUser.ID)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestUserRepository_GetByID_NotFoundError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test"}

	_, err := repo.GetByID(ctx, testUser.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestUserRepository_GetByID_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test"}

	db.Create(&testUser)

	resultUser, err := repo.GetByID(ctx, testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, resultUser.ID, testUser.ID)
}

func TestUserRepository_GetByNickname_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test", Nickname: "Arnold"}
	_ = db.AddError(errors.New("test error"))

	_, err := repo.GetByNickname(ctx, testUser.Nickname)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestUserRepository_GetByNickname_NotFoundError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test", Nickname: "Arnold"}

	_, err := repo.GetByNickname(ctx, testUser.Nickname)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestUserRepository_GetByNickname_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test", Nickname: "Arnold"}

	db.Create(&testUser)

	resultUser, err := repo.GetByNickname(ctx, testUser.Nickname)
	assert.NoError(t, err)
	assert.Equal(t, resultUser.ID, testUser.ID)
}

func TestUserRepository_Update_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test", Nickname: "Arnold"}
	_ = db.AddError(errors.New("test error"))

	_, err := repo.Update(ctx, testUser)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestUserRepository_Update_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase))
	testUser := models.User{ID: "test", Nickname: "Arnold"}

	db.Create(&testUser)
	patchedUser := models.User{ID: testUser.ID, Nickname: "Arnold2"}

	updatedUser, err := repo.Update(ctx, patchedUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.ID, testUser.ID)
	assert.Equal(t, updatedUser.Nickname, patchedUser.Nickname)
}
