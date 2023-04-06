package models

import "gorm.io/gorm"

type Administrator struct {
	gorm.Model
	Email    string `gorm:"not null;unique;index"`
	Password string `gorm:"not null"`
}
