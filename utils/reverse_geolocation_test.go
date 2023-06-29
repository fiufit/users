package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocationFromCoordinates_Ok(t *testing.T) {
	rl, err := NewReverseLocator()
	if err != nil {
		t.Errorf("Error creating ReverseLocator: %v", err)
	}

	lat := 40.416775
	long := -3.703790
	expected := "Spain (ESP), Europe"
	actual, _ := rl.GetLocationFromCoordinates(lat, long)

	assert.Equal(t, expected, actual)
}

func TestGetLocationFromCoordinates_Err(t *testing.T) {
	rl, err := NewReverseLocator()
	if err != nil {
		t.Errorf("Error creating ReverseLocator: %v", err)
	}

	actual, err := rl.GetLocationFromCoordinates(0, 0)

	assert.Equal(t, "", actual)
	assert.Error(t, err)
}
