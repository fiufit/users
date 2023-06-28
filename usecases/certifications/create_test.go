package certifications

import (
	"context"
	"errors"
	"testing"

	"github.com/fiufit/users/contracts"
	certContracts "github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateCertification_ErrUserNotFound(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	creator := NewCertificationCreator(certifications, users)
	ctx := context.Background()

	req := certContracts.CreateCertificationRequest{UserID: "pepe"}

	users.On("GetByID", ctx, req.UserID).Return(models.User{}, contracts.ErrUserNotFound)
	_, err := creator.Create(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestCreateCertification_ErrUserAlreadyCertified(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	creator := NewCertificationCreator(certifications, users)
	ctx := context.Background()

	req := certContracts.CreateCertificationRequest{UserID: "pepe"}

	users.On("GetByID", ctx, req.UserID).Return(models.User{}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusApproved}).Return(certContracts.GetCertificationsResponse{Certifications: []models.Certification{{}}}, nil)

	_, err := creator.Create(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyCertified)
}

func TestCreateCertification_ErrPendingCertsExists(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	creator := NewCertificationCreator(certifications, users)
	ctx := context.Background()

	req := certContracts.CreateCertificationRequest{UserID: "pepe"}

	users.On("GetByID", ctx, req.UserID).Return(models.User{}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusApproved}).Return(certContracts.GetCertificationsResponse{}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusPending}).Return(certContracts.GetCertificationsResponse{Certifications: []models.Certification{{}}}, nil)

	_, err := creator.Create(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrPendingCertsExists)
}

func TestCreateCertification_ErrCreateCertification(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	creator := NewCertificationCreator(certifications, users)
	ctx := context.Background()

	req := certContracts.CreateCertificationRequest{UserID: "pepe"}

	users.On("GetByID", ctx, req.UserID).Return(models.User{}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusApproved}).Return(certContracts.GetCertificationsResponse{}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusPending}).Return(certContracts.GetCertificationsResponse{}, nil)
	certifications.On("Create", ctx, models.Certification{Status: models.CertificationStatusPending, UserID: req.UserID}).Return(models.Certification{}, errors.New("error creating cert"))

	_, err := creator.Create(ctx, req)

	assert.Error(t, err)
}

func TestCreateCertification_Ok(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	creator := NewCertificationCreator(certifications, users)
	ctx := context.Background()

	req := certContracts.CreateCertificationRequest{UserID: "pepe"}

	users.On("GetByID", ctx, req.UserID).Return(models.User{ID: req.UserID}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusApproved}).Return(certContracts.GetCertificationsResponse{}, nil)
	certifications.On("Get", ctx, certContracts.GetCertificationsRequest{UserID: req.UserID, Status: models.CertificationStatusPending}).Return(certContracts.GetCertificationsResponse{}, nil)
	certifications.On("Create", ctx, models.Certification{Status: models.CertificationStatusPending, UserID: req.UserID}).Return(models.Certification{Status: models.CertificationStatusPending}, nil)

	createdCert, err := creator.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, createdCert.User.ID, req.UserID)
	assert.Equal(t, createdCert.Status, models.CertificationStatusPending)
}
