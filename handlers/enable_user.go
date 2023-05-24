package handlers

import (
	"errors"
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

func (h EnableUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		err := h.users.EnableUser(ctx, userID)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(err))
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
				return
			}
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
