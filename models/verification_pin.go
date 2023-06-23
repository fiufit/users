package models

import "time"

type VerificationPin struct {
	UserID    string    `gorm:"primaryKey;not null"`
	Pin       string    `gorm:"not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null"`
}
