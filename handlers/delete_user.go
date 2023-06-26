package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeleteUser struct {
	users  users.UserDeleter
	logger *zap.Logger
}

func NewDeleteUser(users users.UserDeleter, logger *zap.Logger) DeleteUser {
	return DeleteUser{users: users, logger: logger}
}

// User Delete godoc
//	@Summary		Deletes a user by their ID.
//	@Description	Deletes a user by their ID. This endpoint should only be called by admins or the same user. Authorization is the gateway's responsibility.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version						path		string	true	"API Version"
//	@Param			userID						path		string	true	"User ID"
//	@Success		200							{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400							{object}	contracts.ErrResponse
//	@Failure		404							{object}	contracts.ErrResponse
//	@Failure		500							{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID} 	[delete]
func (h DeleteUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		err := h.users.DeleteUser(ctx, userID)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
