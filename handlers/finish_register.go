package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/accounts"
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

// User Register godoc
//	@Summary		Register a new user.
//	@Description	Register a new User. Mandatory to be called after /users/register to complete additional profile info
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string								true	"API Version"
//	@Param			payload								body		ucontracts.FinishRegisterRequest	true	"Body params"
//	@Success		200									{object}	ucontracts.FinishRegisterResponse	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Failure		409									{object}	contracts.ErrResponse
//	@Failure		500									{object}	contracts.ErrResponse
//	@Router			/{version}/users/finish-register 	[post]
func (h FinishRegister) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.FinishRegisterRequest
		err := ctx.ShouldBindJSON(&req)
		validateErr := req.Validate()
		if err != nil || validateErr != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		userID := ctx.MustGet("userID").(string)
		req.UserID = userID

		res, err := h.users.FinishRegister(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
