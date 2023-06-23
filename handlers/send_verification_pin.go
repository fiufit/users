package handlers

import (
	"errors"
	"net/http"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendVerificationPin struct {
	validation accounts.Verificator
	logger     *zap.Logger
}

func NewSendVerificationPin(validation accounts.Verificator, logger *zap.Logger) SendVerificationPin {
	return SendVerificationPin{validation: validation, logger: logger}
}

func (h SendVerificationPin) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.SendVerificationPinRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}
		userID := ctx.MustGet("userID").(string)
		req.UserID = userID
		res, err := h.validation.SendVerificationPin(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(err))
				return
			}

			if errors.Is(err, contracts.ErrUserAlreadyVerified) {
				ctx.JSON(http.StatusConflict, contracts.FormatErrResponse(err))
				return
			}

			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
