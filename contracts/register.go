package contracts

import (
	"time"

	"github.com/fiufit/users/models"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID string `json:"userID"`
}

type FinishRegisterRequest struct {
	UserID       string
	Nickname     string    `json:"nick-name"`
	DisplayName  string    `json:"display-name"`
	IsMale       bool      `json:"is-male"`
	BirthDate    time.Time `json:"birth-date"`
	Height       uint      `json:"height"`
	Weight       uint      `json:"weight"`
	MainLocation string    `json:"main-location"`
	Interests    []string  `json:"-"`
}

type FinishRegisterResponse struct {
	User models.User `json:"user"`
}
