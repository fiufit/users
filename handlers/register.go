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

type Register struct {
	users  accounts.Registerer
	logger *zap.Logger
}

func NewRegister(users accounts.Registerer, logger *zap.Logger) Register {
	return Register{users: users, logger: logger}
}

func (h Register) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.RegisterRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		res, err := h.users.Register(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrUserAlreadyExists) {
				ctx.JSON(http.StatusConflict, contracts.FormatErrResponse(err))
				return
			}

			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
