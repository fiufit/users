package handlers

import (
	"errors"
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUserByNickname struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetUserByNickname(users users.UserGetter, logger *zap.Logger) GetUserByNickname {
	return GetUserByNickname{users: users, logger: logger}
}

func (h GetUserByNickname) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		nickname := ctx.Query("nickname")
		if nickname == "" {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}
		user, err := h.users.GetUserByNickname(ctx, nickname)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
				return
			}
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(user))
	}
}
