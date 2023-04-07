package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Administrator struct {
	gorm.Model
	Email    string `gorm:"not null;unique;index"`
	Password string `gorm:"not null"`
}

func (a Administrator) MarshalJson() ([]byte, error) {
	var tmp struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	}

	tmp.ID = a.ID
	tmp.Email = a.Email

	return json.Marshal(&tmp)
}
