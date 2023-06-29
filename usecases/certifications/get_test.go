package certifications

import (
	"context"
	"errors"
	"testing"

	certContracts "github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCertifications_Error(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	certGetter := NewCertificationGetterImpl(certifications, users)
	ctx := context.Background()

	req := certContracts.GetCertificationsRequest{}

	certifications.On("Get", ctx, req).Return(certContracts.GetCertificationsResponse{}, errors.New("repo error"))

	_, err := certGetter.Get(ctx, req)

	assert.Error(t, err)
}

func TestGetCertifications_Ok(t *testing.T) {
	users := new(mocks.Users)
	certifications := new(mocks.Certifications)
	certGetter := NewCertificationGetterImpl(certifications, users)
	ctx := context.Background()

	req := certContracts.GetCertificationsRequest{}

	users.On("GetByID", ctx, mock.Anything).Return(models.User{}, nil)
	certifications.On("Get", ctx, req).Return(certContracts.GetCertificationsResponse{Certifications: []models.Certification{{}}}, nil)

	_, err := certGetter.Get(ctx, req)

	assert.NoError(t, err)
}
