package handlers

import (
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

// User Register godoc
//
//	@Summary		Register a new user.
//	@Description	Register a new User.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version						path		string						true	"API Version"
//	@Param			payload						body		ucontracts.RegisterRequest	true	"Body params"
//	@Success		200							{object}	ucontracts.RegisterResponse	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400							{object}	contracts.ErrResponse
//	@Failure		409							{object}	contracts.ErrResponse
//	@Failure		500							{object}	contracts.ErrResponse
//	@Router			/{version}/users/register 	[post]
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
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
