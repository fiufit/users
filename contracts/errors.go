package contracts

import "errors"

var (
	ErrInternal            = errors.New("something went wrong")
	ErrBadRequest          = errors.New("unable to parse request")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exist, attempting to re-register same user")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidInterest     = errors.New("invalid interest")
	ErrUserAlreadyDisabled = errors.New("user is already disabled")
	ErrUserNotDisabled     = errors.New("user is not disabled")
	ErrUserAlreadyVerified = errors.New("user is already verified")
)
