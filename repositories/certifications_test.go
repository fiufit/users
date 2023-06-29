package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestCertificationRepository_Create_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)

	repo := NewCertificationRepository(db, logger, firebase)
	db.AddError(errors.New("test db error"))

	_, err := repo.Create(ctx, models.Certification{})
	assert.Error(t, err)
	db.Error = nil
}

func TestCertificationRepository_Create_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)

	repo := NewCertificationRepository(db, logger, firebase)

	testUser := models.User{ID: "testuser"}
	db.Create(&testUser)

	firebase.On("GetCertificationVideoUrl", ctx, testUser.ID).Return("testurl")
	testCert := models.Certification{
		UserID: testUser.ID,
		Status: models.CertificationStatusPending,
	}

	createdCert, err := repo.Create(ctx, testCert)
	assert.NoError(t, err)

	var dbCert models.Certification
	_ = db.First(&dbCert)

	assert.Equal(t, createdCert.UserID, dbCert.UserID)
	assert.Equal(t, createdCert.Status, dbCert.Status)
}

func TestCertificationRepository_GetByID_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)

	repo := NewCertificationRepository(db, logger, firebase)
	db.AddError(errors.New("test db error"))

	_, err := repo.GetByID(ctx, 1)
	assert.Error(t, err)
	db.Error = nil
}

func TestCertificationRepository_GetByID_OK(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)

	repo := NewCertificationRepository(db, logger, firebase)

	testUser := models.User{ID: "testuser"}
	db.Create(&testUser)

	firebase.On("GetCertificationVideoUrl", ctx, testUser.ID).Return("testurl")
	testCert := models.Certification{
		UserID: testUser.ID,
		Status: models.CertificationStatusPending,
	}

	db.Create(&testCert)

	cert, err := repo.GetByID(ctx, testCert.ID)
	assert.NoError(t, err)

	assert.Equal(t, cert.UserID, testCert.UserID)
	assert.Equal(t, cert.Status, testCert.Status)
}

func TestCertificationRepository_Get_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)

	repo := NewCertificationRepository(db, logger, firebase)
	db.AddError(errors.New("test db error"))

	_, err := repo.Get(ctx, certifications.GetCertificationsRequest{})
	assert.Error(t, err)
	db.Error = nil
}

func TestCertificationRepository_Get_Ok(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)
	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)
	firebase.On("GetCertificationVideoUrl", ctx, mock.Anything).Return("testurl")
	repo := NewCertificationRepository(db, logger, firebase)

	testUsers := [2]models.User{
		{ID: "a", Nickname: "a"},
		{ID: "b", Nickname: "b"},
	}

	db.Create(&testUsers)

	testCertifications := [4]models.Certification{
		{UserID: testUsers[0].ID, Status: models.CertificationStatusPending},
		{UserID: testUsers[0].ID, Status: models.CertificationStatusDenied},
		{UserID: testUsers[0].ID, Status: models.CertificationStatusApproved},
		{UserID: testUsers[1].ID, Status: models.CertificationStatusPending},
	}

	db.Create(&testCertifications)

	type testCase struct {
		description       string
		expectedCertCount int
		req               certifications.GetCertificationsRequest
	}

	for _, tCase := range []testCase{
		{
			description:       "Get all certifications",
			expectedCertCount: 4,
			req:               certifications.GetCertificationsRequest{},
		},
		{
			description:       "Get certifications by user",
			expectedCertCount: 3,
			req:               certifications.GetCertificationsRequest{UserID: testUsers[0].ID},
		},
		{
			description:       "Get certifications by status",
			expectedCertCount: 2,
			req:               certifications.GetCertificationsRequest{Status: models.CertificationStatusPending},
		},
		{
			description:       "Get certifications by user and status",
			expectedCertCount: 1,
			req:               certifications.GetCertificationsRequest{UserID: testUsers[0].ID, Status: models.CertificationStatusApproved},
		},
	} {
		t.Run(tCase.description, func(t *testing.T) {
			res, err := repo.Get(ctx, tCase.req)
			assert.NoError(t, err)
			assert.Equal(t, len(res.Certifications), tCase.expectedCertCount)
		})
	}
}

func TestCertificationRepository_Update_DBError(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)
	repo := NewCertificationRepository(db, logger, firebase)
	db.AddError(errors.New("test db error"))

	_, err := repo.Update(ctx, models.Certification{})
	assert.Error(t, err)
	db.Error = nil
}

func TestCertificationRepository_Update_Ok(t *testing.T) {
	defer testSuite.TruncateModels()
	ctx := context.Background()
	db := testSuite.DB.WithContext(ctx)

	logger := zaptest.NewLogger(t)
	firebase := new(mocks.Firebase)
	repo := NewCertificationRepository(db, logger, firebase)

	testUser := models.User{ID: "testuser"}
	db.Create(&testUser)

	firebase.On("GetCertificationVideoUrl", ctx, testUser.ID).Return("testurl")
	testCert := models.Certification{
		UserID: testUser.ID,
		Status: models.CertificationStatusPending,
	}

	db.Create(&testCert)

	testCert.Status = models.CertificationStatusApproved

	updatedCert, err := repo.Update(ctx, testCert)
	assert.NoError(t, err)

	var dbCert models.Certification
	_ = db.First(&dbCert)

	assert.Equal(t, updatedCert.UserID, dbCert.UserID)
	assert.Equal(t, updatedCert.Status, dbCert.Status)
}
