package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
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
		var req contracts.RegisterRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Unable to read request body")
			return
		}

		err = h.users.Register(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "Something went wrong")
			return
		}

		ctx.JSON(http.StatusOK, "OK")
	}
}
