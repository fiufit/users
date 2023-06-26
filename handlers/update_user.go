package handlers

import (
	"errors"
	"net/http"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateUser struct {
	users  users.UserUpdater
	logger *zap.Logger
}

func NewUpdateUser(users users.UserUpdater, logger *zap.Logger) UpdateUser {
	return UpdateUser{users: users, logger: logger}
}

// Update User godoc
//	@Summary		Updates a user.
//	@Description	Updates a user profile info.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version						path		string							true	"API Version"
//	@Param			userID						path		string							true	"User ID"
//	@Param			payload						body		ucontracts.UpdateUserRequest	true	"Body params, all of them optional, ID is ignored and taken from path param"
//	@Success		200							{object}	string							"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400							{object}	contracts.ErrResponse
//	@Failure		404							{object}	contracts.ErrResponse
//	@Failure		409							{object}	contracts.ErrResponse
//	@Failure		500							{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID}	[patch]
func (h UpdateUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.UpdateUserRequest
		err := ctx.ShouldBindJSON(&req)
		validateErr := req.Validate()
		if err != nil || validateErr != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		userID := ctx.MustGet("userID").(string)
		req.ID = userID

		updatedUser, err := h.users.UpdateUser(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
				return
			}
			if errors.Is(err, contracts.ErrUserAlreadyExists) {
				ctx.JSON(http.StatusConflict, contracts.FormatErrResponse(contracts.ErrUserAlreadyExists))
				return
			}
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(updatedUser))
	}
}
