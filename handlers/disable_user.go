package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DisableUser struct {
	users  users.UserEnabler
	logger *zap.Logger
}

func NewDisableUser(users users.UserEnabler, logger *zap.Logger) DisableUser {
	return DisableUser{users: users, logger: logger}
}

// User Disable godoc
//	@Summary		Disables a user by their ID, preventing them from doing further requests.
//	@Description	Disables a user by their ID, preventing them from doing further requests. This endpoint should only be called by admins. Authorization is the gateway's responsibility.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string	true	"API Version"
//	@Param			userID								path		string	true	"User ID"
//	@Success		200									{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Failure		404									{object}	contracts.ErrResponse
//	@Failure		409									{object}	contracts.ErrResponse
//	@Failure		500									{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID}/disable 	[delete]
func (h DisableUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		err := h.users.DisableUser(ctx, userID)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
