package models

import "github.com/fiufit/users/contracts"

var validInterests = map[string]struct{}{
	"strength":    {},
	"speed":       {},
	"endurance":   {},
	"lose weight": {},
	"gain weight": {},
	"sports":      {},
}

type Interest struct {
	Name string `gorm:"primaryKey;not null;index;unique"`
}

func ValidateInterests(interestStrings ...string) ([]Interest, error) {
	interests := make([]Interest, len(interestStrings))
	for i, interest := range interestStrings {
		if _, exists := validInterests[interest]; !exists {
			return []Interest{}, contracts.ErrInvalidInterest
		}
		interests[i] = Interest{Name: interest}
	}
	return interests, nil
}
