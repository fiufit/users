package handlers

import (
	"errors"
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

func (h DeleteUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		err := h.users.DeleteUser(ctx, userID)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
				return
			}
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse("User deleted succesfully"))
	}
}
