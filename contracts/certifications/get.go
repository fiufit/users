package certifications

import (
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
)

type GetCertificationsRequest struct {
	UserID string `form:"user_id"`
	Status string `form:"status"`
	contracts.Pagination
}

type GetCertificationsResponse struct {
	Certifications []models.Certification `json:"certifications"`
	contracts.Pagination
}
