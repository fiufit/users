package handlers

import (
	"errors"
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FinishRegister struct {
	users  accounts.Registerer
	logger *zap.Logger
}

func NewFinishRegister(users accounts.Registerer, logger *zap.Logger) FinishRegister {
	return FinishRegister{users: users, logger: logger}
}

func (h FinishRegister) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req contracts.FinishRegisterRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Unable to read request body")
			return
		}

		err = h.users.FinishRegister(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrUserAlreadyExists) {
				ctx.JSON(http.StatusConflict, err.Error())
				return
			}
			ctx.JSON(http.StatusInternalServerError, "Something went wrong")
			return
		}

		ctx.JSON(http.StatusOK, "OK")
	}
}