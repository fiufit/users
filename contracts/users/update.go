package users

import (
	"time"

	"github.com/fiufit/users/models"
)

type UpdateUserRequest struct {
	ID           string
	Nickname     string    `json:"nickname" `
	DisplayName  string    `json:"display_name" `
	IsMale       bool      `json:"is_male" `
	BirthDate    time.Time `json:"birth_date" `
	Height       uint      `json:"height" `
	Weight       uint      `json:"weight" `
	MainLocation string    `json:"main_location" `
}

type UpdateUserResponse struct {
	User models.User `json:"user"`
}
