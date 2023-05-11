package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

func TestAdminRepository_Create_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()

	db := testSuite.DB
	repository := NewAdminRepository(db, zaptest.NewLogger(t))

	_ = repository.db.AddError(errors.New("test error"))
	testAdmin := models.Administrator{}

	_, err := repository.Create(ctx, testAdmin)
	assert.Error(t, err)
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestAdminRepository_Create_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()

	db := testSuite.DB
	repository := NewAdminRepository(db, zaptest.NewLogger(t))

	testAdmin := models.Administrator{
		Model:    gorm.Model{},
		Email:    "testadmin@testemail.com",
		Password: "testtest",
	}
	createdAdmin, err := repository.Create(ctx, testAdmin)

	assert.NoError(t, err)

	var dbAdmin models.Administrator
	_ = db.First(&dbAdmin)

	assert.Equal(t, createdAdmin.Email, dbAdmin.Email)
}

func TestAdminRepository_Create_DuplicatedEmailError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()

	db := testSuite.DB
	repository := NewAdminRepository(db, zaptest.NewLogger(t))

	testAdmin := models.Administrator{
		Model:    gorm.Model{},
		Email:    "testadmin@testemail.com",
		Password: "testtest",
	}

	_ = db.Create(&testAdmin)

	_, err := repository.Create(ctx, testAdmin)
	assert.Error(t, err)
	assert.Equal(t, contracts.ErrUserAlreadyExists, err)

	db.Error = nil
}

func TestAdminRepository_GetByEmail_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()

	db := testSuite.DB
	repository := NewAdminRepository(db, zaptest.NewLogger(t))

	_ = repository.db.AddError(gorm.ErrUnsupportedDriver)

	_, err := repository.GetByEmail(ctx, "testadmin@fiufit.com")
	assert.Error(t, err)
	db.Error = nil //overwrite the db Error so that TruncateModels() doesn't panic
}

func TestAdminRepository_GetByEmail_NotFound(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()

	db := testSuite.DB
	repository := NewAdminRepository(db, zaptest.NewLogger(t))

	_, err := repository.GetByEmail(ctx, "testadmin@fiufit.com")
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestAdminRepository_GetByEmail_Ok(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()

	db := testSuite.DB
	repository := NewAdminRepository(db, zaptest.NewLogger(t))

	testAdmin := models.Administrator{
		Model:    gorm.Model{},
		Email:    "testadmin@fiufit.com",
		Password: "testtest",
	}

	_ = db.Create(&testAdmin)

	resultAdmin, err := repository.GetByEmail(ctx, "testadmin@fiufit.com")
	assert.NoError(t, err)
	assert.Equal(t, resultAdmin.Email, testAdmin.Email)
}
