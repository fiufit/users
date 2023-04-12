package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string    `gorm:"primaryKey;not null"`
	Nickname          string    `gorm:"not null;unique;index"`
	DisplayName       string    `gorm:"not null"`
	IsMale            bool      `gorm:"not null"`
	CreatedAt         time.Time `gorm:"not null"`
	DeletedAt         gorm.DeletedAt
	BornAt            time.Time `gorm:"not null"`
	Height            uint      `gorm:"not null"`
	Weight            uint      `gorm:"not null"`
	IsVerifiedTrainer bool      `gorm:"not null;default:false"`
	MainLocation      string    `gorm:"not null"`
	Interests         []string  `gorm:"-"`
}

func (u User) ToPublicView() map[string]interface{} {
	return map[string]interface{}{
		"id":            u.ID,
		"nickname":      u.Nickname,
		"display_name":  u.DisplayName,
		"is_male":       u.IsMale,
		"is_verified":   u.IsVerifiedTrainer,
		"main_location": u.MainLocation,
	}
}

func (u User) ToPrivilegedView() map[string]interface{} {
	userMap := u.ToPublicView()
	userMap["creation_date"] = u.CreatedAt
	userMap["birth_date"] = u.BornAt
	userMap["height"] = u.Height
	userMap["weight"] = u.Weight

	return userMap
}
