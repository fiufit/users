package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiufit/users/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestVerificationPinRepository_Create_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)
	logger := zaptest.NewLogger(t)

	repo := NewVerificationPinRepository(db, logger)
	db.AddError(errors.New("test db error"))

	_, err := repo.Create(ctx, models.VerificationPin{})
	assert.Error(t, err)
	db.Error = nil
}

func TestVerificationPinRepository_Create_Ok(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)
	logger := zaptest.NewLogger(t)
	testUser := models.User{ID: "a", Nickname: "a"}
	testPin := models.VerificationPin{
		UserID:    "a",
		Pin:       "",
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}
	db.Create(&testUser)

	repo := NewVerificationPinRepository(db, logger)
	repo.Create(ctx, testPin)

	var dbPin models.VerificationPin
	_ = db.First(&dbPin)

	assert.Equal(t, testPin.UserID, dbPin.UserID)
	assert.Equal(t, testPin.Pin, dbPin.Pin)
}

func TestVerificationPinRepository_GetByUserID_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)
	logger := zaptest.NewLogger(t)

	repo := NewVerificationPinRepository(db, logger)
	db.AddError(errors.New("test db error"))

	_, err := repo.GetByUserID(ctx, "a")
	assert.Error(t, err)
	db.Error = nil
}

func TestVerificationPinRepository_GetByUserID_ErrNotFound(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)
	logger := zaptest.NewLogger(t)

	repo := NewVerificationPinRepository(db, logger)

	_, err := repo.GetByUserID(ctx, "a")
	assert.Error(t, err)
}

func TestVerificationPinRepository_GetByUserID_Ok(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)
	logger := zaptest.NewLogger(t)
	testUser := models.User{ID: "a", Nickname: "a"}
	testPin := models.VerificationPin{
		UserID:    "a",
		Pin:       "",
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}
	db.Create(&testUser)
	db.Create(&testPin)

	repo := NewVerificationPinRepository(db, logger)
	pin, err := repo.GetByUserID(ctx, "a")
	assert.NoError(t, err)
	assert.Equal(t, testPin.UserID, pin.UserID)
	assert.Equal(t, testPin.Pin, pin.Pin)
}
