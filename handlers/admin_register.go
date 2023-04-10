package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	accounts2 "github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminRegister struct {
	admins accounts.AdminRegisterer
	logger *zap.Logger
}

func NewAdminRegister(admins accounts.AdminRegisterer, logger *zap.Logger) AdminRegister {
	return AdminRegister{admins: admins, logger: logger}
}

func (h AdminRegister) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req accounts2.AdminRegisterRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		res, err := h.admins.Register(ctx, req)
		if err != nil {
			h.logger.Error("unable to register new administrator", zap.Error(err), zap.Any("request", req))
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
