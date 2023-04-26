package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordErr(t *testing.T) {
	assert.Error(t, ValidatePassword("hunter", "hunter2"))
}

func TestValidatePasswordOk(t *testing.T) {
	assert.Error(t, ValidatePassword("hunter2", "hunter2"))
}

func TestHashPassWordOk(t *testing.T) {
	//hunter2 with 10 rounds of Bcrypt
	password := "hunter2"
	hashedPass, _ := HashPassword(password)

	err := ValidatePassword(password, hashedPass)
	assert.NoError(t, err)
}

func TestHashPasswordTooLongError(t *testing.T) {
	password := "huuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuunter2"
	_, err := HashPassword(password)
	assert.Error(t, err)
}
