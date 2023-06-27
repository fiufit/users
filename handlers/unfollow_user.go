package handlers

import (
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

// User Unfollow godoc
//
//	@Summary		Unfollow an user.
//	@Description	Removes a following relationship between two users.
//	@Tags			followers
//	@Accept			json
//	@Produce		json
//	@Param			version												path		string	true	"API Version"
//	@Param			followerID											path		string	true	"userID of the following user"
//	@Param			userID												path		string	true	"userID of followed user"
//	@Success		200													{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400													{object}	contracts.ErrResponse
//	@Failure		404													{object}	contracts.ErrResponse
//	@Failure		500													{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID}/followers/{followerID} 	[delete]
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
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
