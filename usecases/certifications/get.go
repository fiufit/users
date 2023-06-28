package certifications

import (
	"context"

	"github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/repositories"
)

type CertificationGetter interface {
	Get(ctx context.Context, req certifications.GetCertificationsRequest) (certifications.GetCertificationsResponse, error)
}

type CertificationGetterImpl struct {
	certifications repositories.Certifications
	users          repositories.Users
}

func NewCertificationGetterImpl(certifications repositories.Certifications, users repositories.Users) CertificationGetterImpl {
	return CertificationGetterImpl{certifications: certifications, users: users}
}

func (uc CertificationGetterImpl) Get(ctx context.Context, req certifications.GetCertificationsRequest) (certifications.GetCertificationsResponse, error) {
	res, err := uc.certifications.Get(ctx, req)
	if err != nil {
		return certifications.GetCertificationsResponse{}, err
	}

	for i, cert := range res.Certifications {
		certUser, _ := uc.users.GetByID(ctx, cert.UserID)
		res.Certifications[i].User = certUser
	}
	return res, nil
}
