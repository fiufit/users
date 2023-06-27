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
}

func NewCertificationGetterImpl(certifications repositories.Certifications) CertificationGetterImpl {
	return CertificationGetterImpl{certifications: certifications}
}

func (uc CertificationGetterImpl) Get(ctx context.Context, req certifications.GetCertificationsRequest) (certifications.GetCertificationsResponse, error) {
	return uc.certifications.Get(ctx, req)
}
