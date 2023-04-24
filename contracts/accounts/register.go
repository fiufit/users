package accounts

import (
	"time"

	"github.com/fiufit/users/models"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	UserID string `json:"userID"`
}

type FinishRegisterRequest struct {
	UserID       string
	Nickname     string    `json:"nickname" binding:"required"`
	DisplayName  string    `json:"display_name" binding:"required"`
	IsMale       *bool     `json:"is_male" binding:"required"`
	BirthDate    time.Time `json:"birth_date" binding:"required"`
	Height       uint      `json:"height" binding:"required"`
	Weight       uint      `json:"weight" binding:"required"`
	MainLocation string    `json:"main_location" binding:"required"`
	Interests    []string  `json:"-"`
}

type FinishRegisterResponse struct {
	User models.User `json:"user"`
}
