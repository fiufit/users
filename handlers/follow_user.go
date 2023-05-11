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

type FollowUser struct {
	follows users.UserFollower
	logger  *zap.Logger
}

func NewFollowUser(follows users.UserFollower, logger *zap.Logger) FollowUser {
	return FollowUser{follows: follows, logger: logger}
}

func (h FollowUser) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.FollowUserRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		req.FollowedUserID = ctx.MustGet("userID").(string)

		err = h.follows.FollowUser(ctx, req)
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
