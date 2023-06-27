package certifications

import (
	"context"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
)

type CertificationUpdater interface {
	Update(ctx context.Context, req certifications.UpdateCertificationRequest) (models.Certification, error)
}

type CertificationUpdaterImpl struct {
	certifications repositories.Certifications
	users          repositories.Users
	notifications  repositories.Notifications
	firebase       repositories.Firebase
	logger         *zap.Logger
}

func NewCertificationUpdaterImpl(certifications repositories.Certifications, users repositories.Users, notifications repositories.Notifications, firebase repositories.Firebase, logger *zap.Logger) CertificationUpdaterImpl {
	return CertificationUpdaterImpl{certifications: certifications, users: users, notifications: notifications, firebase: firebase, logger: logger}
}

func (uc CertificationUpdaterImpl) Update(ctx context.Context, req certifications.UpdateCertificationRequest) (models.Certification, error) {
	cert, err := uc.certifications.GetByID(ctx, req.CertificationID)
	if err != nil {
		return models.Certification{}, err
	}

	if cert.Status == models.CertificationStatusApproved {
		if req.Status == models.CertificationStatusApproved {
			return cert, nil
		}
		return models.Certification{}, contracts.ErrUserAlreadyCertified
	}

	user, err := uc.users.GetByID(ctx, cert.UserID)
	if err != nil {
		return models.Certification{}, err
	}

	cert.Status = req.Status
	updatedCert, err := uc.certifications.Update(ctx, cert)
	if err != nil {
		return models.Certification{}, err
	}

	if updatedCert.Status == models.CertificationStatusApproved {
		user.IsVerifiedTrainer = true
		_, err := uc.users.Update(ctx, user)
		if err != nil {
			return models.Certification{}, err
		}
	}

	if updatedCert.Status != models.CertificationStatusPending {
		err = uc.notifications.SendCertificationNotification(ctx, user.ID, updatedCert.Status)
		if err != nil {
			uc.logger.Error("Unable to send Certification status update notification", zap.Error(err), zap.Any("cert", updatedCert))
		}
	}

	updatedCert.User = user
	updatedCert.VideoUrl = uc.firebase.GetCertificationVideoUrl(ctx, user.ID)
	return updatedCert, nil
}
