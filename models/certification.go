package models

import "gorm.io/gorm"

const CertificationStatusPending = "pending"
const CertificationStatusDenied = "denied"
const CertificationStatusAccepted = "accepted"

var validCertificationStatuses = map[string]struct{}{
	CertificationStatusPending:  {},
	CertificationStatusDenied:   {},
	CertificationStatusAccepted: {},
}

type Certification struct {
	gorm.Model
	UserID   string
	User     User
	Status   string
	VideoUrl string
}
