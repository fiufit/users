package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/accounts"
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

// Admin Register godoc
//	@Summary		Register an administrator
//	@Description	Register a new admin. This endpoint should only be called after a gateway processed the corresponding authorization
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version						path		string							true	"API Version"
//	@Param			payload						body		ucontracts.AdminRegisterRequest	true	"Body params"
//	@Success		200							{object}	accounts.AdminRegisterResponse	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400							{object}	contracts.ErrResponse
//	@Failure		409							{object}	contracts.ErrResponse
//	@Failure		500							{object}	contracts.ErrResponse
//	@Router			/{version}/admin/register 	[post]
func (h AdminRegister) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.AdminRegisterRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		res, err := h.admins.Register(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
