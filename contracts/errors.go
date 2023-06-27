package contracts

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInternal               = errors.New("something went wrong")
	ErrBadRequest             = errors.New("unable to parse request")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exist, attempting to re-register same user")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrInvalidInterest        = errors.New("invalid interest")
	ErrUserAlreadyDisabled    = errors.New("user is already disabled")
	ErrUserNotDisabled        = errors.New("user is not disabled")
	ErrUserAlreadyVerified    = errors.New("user is already verified")
	ErrVerificationPinExpired = errors.New("verification pin is expired")
	ErrInvalidVerificationPin = errors.New("invalid verification pin")
	ErrUserAlreadyCertified   = errors.New("user is already verified")
	ErrPendingCertsExists     = errors.New("A pending certification request already exists")
	ErrCertificationNotFound  = errors.New("certification not found")
)

func HandleErrorType(ctx *gin.Context, err error) {
	var status int

	switch {
	case errors.Is(err, ErrBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, ErrUserNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrUserAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, ErrInvalidPassword):
		status = http.StatusUnauthorized
	case errors.Is(err, ErrInvalidInterest):
		status = http.StatusBadRequest
	case errors.Is(err, ErrUserAlreadyDisabled):
		status = http.StatusConflict
	case errors.Is(err, ErrUserNotDisabled):
		status = http.StatusConflict
	case errors.Is(err, ErrUserAlreadyVerified):
		status = http.StatusConflict
	case errors.Is(err, ErrUserAlreadyCertified):
		status = http.StatusConflict
	case errors.Is(err, ErrPendingCertsExists):
		status = http.StatusConflict
	case errors.Is(err, ErrCertificationNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrVerificationPinExpired):
		status = http.StatusUnauthorized
	case errors.Is(err, ErrInvalidVerificationPin):
		status = http.StatusUnauthorized
	default:
		status = http.StatusInternalServerError
		ctx.JSON(status, FormatErrResponse(ErrInternal))
		return
	}

	ctx.JSON(status, FormatErrResponse(err))
}
