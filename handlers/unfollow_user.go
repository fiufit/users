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

type UnfollowUser struct {
	follows users.UserFollower
	logger  *zap.Logger
}

func NewUnfollowUser(follows users.UserFollower, logger *zap.Logger) UnfollowUser {
	return UnfollowUser{follows: follows, logger: logger}
}

func (h UnfollowUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.UnfollowUserRequest
		err := ctx.ShouldBindUri(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		req.FollowedUserID = ctx.MustGet("userID").(string)

		err = h.follows.UnfollowUser(ctx, req)
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
