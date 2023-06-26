package handlers

import (
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

// User Follow godoc
//	@Summary		Follow an user.
//	@Description	Creates a following relationship from the requesting user to the one in the route.
//	@Tags			followers
//	@Accept			json
//	@Produce		json
//	@Param			version									path		string	true	"API Version"
//	@Param			follower_id								query		string	true	"userID of the following user"
//	@Param			userID									path		string	true	"userID of followed user"
//	@Success		200										{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400										{object}	contracts.ErrResponse
//	@Failure		404										{object}	contracts.ErrResponse
//	@Failure		500										{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID}/followers 	[post]
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
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
