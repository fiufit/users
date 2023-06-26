package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	uContracts "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetClosestUsers struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetClosestUsers(users users.UserGetter, logger *zap.Logger) GetClosestUsers {
	return GetClosestUsers{users: users, logger: logger}
}

// Get Closest Users godoc
//
//	@Summary		Gets the closest users to a central user.
//	@Description	Gets the closest users to a central user.
//	@Tags			followers
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string					true	"API Version"
//	@Param			userID								path		string					true	"userID of the person whose near users we want to find"
//	@Param			distance							query		int						true	"distance radio (meters) in which to find users"
//	@Param			page								query		int						false	"page number when getting with pagination"
//	@Param			page_size							query		int						false	"page size when getting with pagination"
//	@Success		200									{object}	users.GetUsersResponse	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Failure		404									{object}	contracts.ErrResponse
//	@Failure		500									{object}	contracts.ErrResponse
//	@Router			/{version}/users/{userID}/closest	[get]
func (h GetClosestUsers) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req uContracts.GetClosestUsersRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		req.UserID = ctx.MustGet("userID").(string)

		req.Pagination.Validate()
		resUsers, err := h.users.GetClosestUsers(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(resUsers))
	}
}
