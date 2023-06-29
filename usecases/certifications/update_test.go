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
	"go.uber.org/zap/zaptest"
)

func TestUpdateCertifications_GetCertError(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "approved",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{}, contracts.ErrCertificationNotFound)
	_, err := certUpdater.Update(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrCertificationNotFound)
}

func TestUpdateCertifications_ErrUserAlreadyCertified(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "denied",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{Status: models.CertificationStatusApproved}, nil)
	_, err := certUpdater.Update(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyCertified)
}

func TestUpdateCertifications_UserAlreadyCertifiedOk(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "approved",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{Status: models.CertificationStatusApproved}, nil)
	_, err := certUpdater.Update(ctx, req)

	assert.NoError(t, err)
}

func TestUpdateCertifications_GetUserError(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "approved",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{UserID: "a"}, nil)
	users.On("GetByID", ctx, "a").Return(models.User{}, contracts.ErrUserNotFound)
	_, err := certUpdater.Update(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestUpdateCertifications_UpdateCertificationError(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "approved",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{UserID: "a"}, nil)
	users.On("GetByID", ctx, "a").Return(models.User{ID: "a"}, nil)
	certifications.On("Update", ctx, models.Certification{UserID: "a", Status: req.Status}).Return(models.Certification{}, errors.New("repo update error"))
	_, err := certUpdater.Update(ctx, req)

	assert.Error(t, err)
}

func TestUpdateCertification_UpdateUserError(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "approved",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{UserID: "a"}, nil)
	users.On("GetByID", ctx, "a").Return(models.User{ID: "a"}, nil)
	certifications.On("Update", ctx, models.Certification{UserID: "a", Status: req.Status}).Return(models.Certification{UserID: "a", Status: req.Status}, nil)
	users.On("Update", ctx, models.User{ID: "a", IsVerifiedTrainer: true}).Return(models.User{}, errors.New("repo update error"))
	_, err := certUpdater.Update(ctx, req)

	assert.Error(t, err)
}

func TestUpdateCertification_Ok(t *testing.T) {
	certifications := new(mocks.Certifications)
	users := new(mocks.Users)
	notifications := new(mocks.Notifications)
	firebase := new(mocks.Firebase)
	logger := zaptest.NewLogger(t)
	certUpdater := NewCertificationUpdaterImpl(certifications, users, notifications, firebase, logger)
	ctx := context.Background()

	req := certContracts.UpdateCertificationRequest{
		CertificationID: 1,
		Status:          "approved",
	}

	certifications.On("GetByID", ctx, req.CertificationID).Return(models.Certification{UserID: "a"}, nil)
	users.On("GetByID", ctx, "a").Return(models.User{ID: "a"}, nil)
	certifications.On("Update", ctx, models.Certification{UserID: "a", Status: req.Status}).Return(models.Certification{UserID: "a", Status: req.Status}, nil)
	users.On("Update", ctx, models.User{ID: "a", IsVerifiedTrainer: true}).Return(models.User{}, nil)
	notifications.On("SendCertificationNotification", ctx, "a", req.Status).Return(nil)
	firebase.On("GetCertificationVideoUrl", ctx, "a").Return("url", nil)
	_, err := certUpdater.Update(ctx, req)

	assert.NoError(t, err)
}
