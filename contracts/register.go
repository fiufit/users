package contracts

import (
	"time"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
