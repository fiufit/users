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

type GetFollowedUsers struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetFollowedUsers(users users.UserGetter, logger *zap.Logger) GetFollowedUsers {
	return GetFollowedUsers{users: users, logger: logger}
}

func (h GetFollowedUsers) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.GetFollowedUsersRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		req.UserID = ctx.MustGet("userID").(string)

		res, err := h.users.GetUserFollowed(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
				return
			}
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
