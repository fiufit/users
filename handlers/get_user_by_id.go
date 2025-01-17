package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUserByID struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetUserByID(users users.UserGetter, logger *zap.Logger) GetUserByID {
	return GetUserByID{users: users, logger: logger}
}

// Get User by ID godoc
//
//	@Summary		Gets a user by their ID.
//	@Description	Gets a user by their ID.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version						path		string		true	"API Version"
//	@Param			userID						path		string		true	"User ID"
//	@Success		200							{object}	models.User	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400							{object}	contracts.ErrResponse
//	@Failure		404							{object}	contracts.ErrResponse
//	@Failure		500							{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID} 	[get]
func (h GetUserByID) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		user, err := h.users.GetUserByID(ctx, userID)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(user))
	}
}
