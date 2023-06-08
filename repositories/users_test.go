package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/fiufit/users/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

func TestUserRepository_CreateUser_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
	testUser := models.User{ID: "test"}

	_, err := repo.GetByID(ctx, testUser.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestUserRepository_GetByID_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
	testUser := models.User{ID: "test", Nickname: "Arnold"}

	_, err := repo.GetByNickname(ctx, testUser.Nickname)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestUserRepository_GetByNickname_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
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
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)
	testUser := models.User{ID: "test", Nickname: "Arnold"}

	db.Create(&testUser)
	patchedUser := models.User{ID: testUser.ID, Nickname: "Arnold2"}

	updatedUser, err := repo.Update(ctx, patchedUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.ID, testUser.ID)
	assert.Equal(t, updatedUser.Nickname, patchedUser.Nickname)
}

func TestUserRepository_DeleteUser_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	fireBaseMock := new(mocks.Firebase)
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), fireBaseMock, reverseLocator)
	testUser := models.User{ID: "testUserID"}

	_ = db.Create(&testUser)
	_ = db.AddError(errors.New("test error"))
	err := repo.DeleteUser(ctx, testUser.ID)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic

	var existingUser models.User
	res := db.Where("id = ?", testUser.ID).First(&existingUser)
	assert.NoError(t, res.Error)
	assert.Equal(t, existingUser.ID, testUser.ID)
}

func TestUserRepository_DeleteUser_FirebaseError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	fireBaseMock := new(mocks.Firebase)
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), fireBaseMock, reverseLocator)
	testUser := models.User{ID: "testUserID"}

	fireBaseMock.On("DeleteUser", ctx, testUser.ID).Return(errors.New("test error"))
	_ = db.Create(&testUser)
	err := repo.DeleteUser(ctx, testUser.ID)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic

	var existingUser models.User
	res := db.Where("id = ?", testUser.ID).First(&existingUser)
	assert.NoError(t, res.Error)
	assert.Equal(t, existingUser.ID, testUser.ID)
}

func TestUserRepository_DeleteUser_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	fireBaseMock := new(mocks.Firebase)
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), fireBaseMock, reverseLocator)
	testUser := models.User{ID: "testID"}

	_ = db.Create(&testUser)
	fireBaseMock.On("DeleteUser", ctx, testUser.ID).Return(nil)
	err := repo.DeleteUser(ctx, testUser.ID)

	assert.NoError(t, err)
	var existingUser models.User
	result := db.First(&existingUser)
	assert.Error(t, result.Error)
	assert.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
}

func TestUserRepository_Get_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)

	_ = db.AddError(errors.New("test error"))
	testReq := users.GetUsersRequest{}
	_, err := repo.Get(ctx, testReq)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestUserRepository_Get_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB
	reverseLocator, _ := utils.NewReverseLocator()
	repo := NewUserRepository(db, zaptest.NewLogger(t), new(mocks.Firebase), reverseLocator)

	testUsers := [7]models.User{
		models.User{ID: "a", Nickname: "Guille"},
		models.User{ID: "b", Nickname: "Goye", DisplayName: "Arnie"},
		models.User{ID: "c", Nickname: "Guillermo"},
		models.User{ID: "d", Nickname: "Bob"},
		models.User{ID: "e", Nickname: "Arnold", IsVerifiedTrainer: true},
		models.User{ID: "f", Nickname: "Arnolda"},
		models.User{ID: "g", Nickname: "Mike", IsVerifiedTrainer: true},
	}

	_ = db.Create(&testUsers)

	type testCase struct {
		description string
		expectedIDs []string
		req         users.GetUsersRequest
	}

	auxTrue := true   //lol, we have to do this to mock a request with a pointer to bool
	auxFalse := false //lol, we have to do this to mock a request with a pointer to bool

	for _, tcase := range []testCase{
		{
			description: "GetWithVerifiedTrue",
			expectedIDs: []string{"e", "g"},
			req:         users.GetUsersRequest{IsVerified: &auxTrue},
		},
		{
			description: "GetWithVerifiedFalse",
			expectedIDs: []string{"a", "b", "c", "d", "f"},
			req:         users.GetUsersRequest{IsVerified: &auxFalse},
		},
		{
			description: "GetWithNameLikeArn",
			expectedIDs: []string{"b", "e", "f"},
			req:         users.GetUsersRequest{Name: "Arn"},
		},
		{
			description: "GetAll",
			expectedIDs: []string{"a", "b", "c", "d", "e", "f", "g"},
			req:         users.GetUsersRequest{},
		},
	} {
		t.Run(tcase.description, func(t *testing.T) {
			resultUsers, err := repo.Get(ctx, tcase.req)
			assert.NoError(t, err)
			assert.True(t, areUserIDsInResult(tcase.expectedIDs, resultUsers.Users))
		})
	}

}

func areUserIDsInResult(ids []string, users []models.User) bool {
	isIncluded := false

	for _, id := range ids {
		for _, user := range users {
			if user.ID == id {
				isIncluded = true
			}
		}
		if !isIncluded {
			return false
		}
	}
	return true
}
