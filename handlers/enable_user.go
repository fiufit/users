package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EnableUser struct {
	users  users.UserEnabler
	logger *zap.Logger
}

func NewEnableUser(users users.UserEnabler, logger *zap.Logger) EnableUser {
	return EnableUser{users: users, logger: logger}
}

// User Enable godoc
//	@Summary		Re-enables a user by their ID, allowing them to do further requests.
//	@Description	Re-enables a user by their ID, allowing them to do further requests. This endpoint should only be called by admins. Authorization is the gateway's responsibility.
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
//	@Router			/{version}/users/{userID}/enable 	[post]
func (h EnableUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		err := h.users.EnableUser(ctx, userID)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
