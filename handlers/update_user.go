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

func (h UpdateUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.UpdateUserRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
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
