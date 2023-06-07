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
	UserID          string
	Nickname        string            `json:"nickname" binding:"required"`
	DisplayName     string            `json:"display_name" binding:"required"`
	IsMale          *bool             `json:"is_male" binding:"required"`
	BirthDate       time.Time         `json:"birth_date" binding:"required"`
	Height          uint              `json:"height" binding:"required"`
	Weight          uint              `json:"weight" binding:"required"`
	MainLocation    string            `json:"main_location" binding:"required"`
	InterestStrings []string          `json:"interests"`
	Interests       []models.Interest `json:"-"`
	Method          string            `json:"method"`
}

func (req *FinishRegisterRequest) Validate() error {
	interests, err := models.ValidateInterests(req.InterestStrings...)
	if err != nil {
		return err
	}
	err = ValidateMethod(req.Method)
	if err != nil {
		return err
	}

	req.Interests = interests
	return nil
}

type FinishRegisterResponse struct {
	User models.User `json:"user"`
}
