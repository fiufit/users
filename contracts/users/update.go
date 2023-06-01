package users

import (
	"time"

	"github.com/fiufit/users/models"
)

type UpdateUserRequest struct {
	ID              string
	Nickname        string            `json:"nickname" `
	DisplayName     string            `json:"display_name" `
	IsMale          *bool             `json:"is_male" `
	BirthDate       time.Time         `json:"birth_date" `
	Height          uint              `json:"height" `
	Weight          uint              `json:"weight" `
	Latitude        float64           `json:"latitude"`
	Longitude       float64           `json:"longitude"`
	InterestStrings []string          `json:"interests"`
	Interests       []models.Interest `json:"-"`
}

func (req *UpdateUserRequest) Validate() error {
	interests, err := models.ValidateInterests(req.InterestStrings...)
	if err != nil {
		return err
	}
	req.Interests = interests
	return nil
}

type UpdateUserResponse struct {
	User models.User `json:"user"`
}
