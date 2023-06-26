package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUserFollowers struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetUserFollowers(users users.UserGetter, logger *zap.Logger) GetUserFollowers {
	return GetUserFollowers{users: users, logger: logger}
}

// Get User Followers godoc
//
//	@Summary		Gets the followers of a user.
//	@Description	Gets the followers of a user.
//	@Tags			followers
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string							true	"API Version"
//	@Param			userID								path		string							true	"userID of the person whose followers we want to GET"
//	@Param			page								query		int								false	"page number when getting with pagination"
//	@Param			page_size							query		int								false	"page size when getting with pagination"
//	@Success		200									{object}	users.GetUserFollowersResponse	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Failure		404									{object}	contracts.ErrResponse
//	@Failure		500									{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID}/followers	[get]
func (h GetUserFollowers) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ucontracts.GetUserFollowersRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		req.UserID = ctx.MustGet("userID").(string)

		res, err := h.users.GetUserFollowers(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
