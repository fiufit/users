package models

import "gorm.io/gorm"

const CertificationStatusPending = "pending"
const CertificationStatusDenied = "denied"
const CertificationStatusApproved = "approved"

var validCertificationStatuses = map[string]struct{}{
	CertificationStatusPending:  {},
	CertificationStatusDenied:   {},
	CertificationStatusApproved: {},
}

type Certification struct {
	gorm.Model
	UserID   string
	User     User
	Status   string
	VideoUrl string `gorm:"-"`
}
