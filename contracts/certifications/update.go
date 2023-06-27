package certifications

type UpdateCertificationRequest struct {
	CertificationID uint
	Status          string `form:"status"`
}
