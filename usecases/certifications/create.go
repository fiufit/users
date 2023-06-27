package certifications

import (
	"context"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
)

type CertificationCreator interface {
	Create(ctx context.Context, request certifications.CreateCertificationRequest) (models.Certification, error)
}

type CertificationCreatorImpl struct {
	users          repositories.Users
	certifications repositories.Certifications
}

func NewCertificationCreator(certifications repositories.Certifications, users repositories.Users) CertificationCreatorImpl {
	return CertificationCreatorImpl{certifications: certifications, users: users}
}

func (uc CertificationCreatorImpl) Create(ctx context.Context, request certifications.CreateCertificationRequest) (models.Certification, error) {
	_, err := uc.users.GetByID(ctx, request.UserID)
	if err != nil {
		return models.Certification{}, err
	}

	approvedCerts, err := uc.certifications.Get(ctx, certifications.GetCertificationsRequest{UserID: request.UserID, Status: models.CertificationStatusApproved})
	if err != nil {
		return models.Certification{}, err
	}

	if len(approvedCerts.Certifications) > 0 {
		return models.Certification{}, contracts.ErrUserAlreadyCertified
	}

	pendingCerts, err := uc.certifications.Get(ctx, certifications.GetCertificationsRequest{UserID: request.UserID, Status: models.CertificationStatusPending})
	if err != nil {
		return models.Certification{}, err
	}

	if len(pendingCerts.Certifications) > 0 {
		return models.Certification{}, contracts.ErrPendingCertsExists
	}

	cert := models.Certification{Status: models.CertificationStatusPending, UserID: request.UserID}
	return uc.certifications.Create(ctx, cert)
}
