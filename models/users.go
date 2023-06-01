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
	BornAt            time.Time  `gorm:"not null"`
	Height            uint       `gorm:"not null"`
	Weight            uint       `gorm:"not null"`
	IsVerifiedTrainer bool       `gorm:"not null;default:false"`
	Followers         []User     `gorm:"many2many:user_followers"`
	Latitude          float64    `gorm:"not null"`
	Longitude         float64    `gorm:"not null"`
	Interests         []Interest `gorm:"many2many:user_interests"`
	Disabled          bool       `gorm:"not null"`
	PictureUrl        string     `gorm:"-"`
}

func (u User) ToPublicView() map[string]interface{} {
	return map[string]interface{}{
		"id":           u.ID,
		"nickname":     u.Nickname,
		"display_name": u.DisplayName,
		"is_male":      u.IsMale,
		"is_verified":  u.IsVerifiedTrainer,
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
