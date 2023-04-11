package handlers

import (
	"errors"
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

func (h GetUserByID) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		user, err := h.users.GetUserByID(ctx, userID)
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
