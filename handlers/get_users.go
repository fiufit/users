package handlers

import (
	"errors"
	"net/http"

	"github.com/fiufit/users/contracts"
	users2 "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/usecases/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUsers struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetUsers(users users.UserGetter, logger *zap.Logger) GetUsers {
	return GetUsers{users: users, logger: logger}
}

// Get Users godoc
// @Summary      Gets users by different query params with pagination.
// @Description	 Gets users by their name, nickname, location or verification status. If nickname has a value, other parameters are ignored.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        version   path      string  true  "API Version"
// @Param        nickname   query      string  false  "User Nickname"
// @Param        name   query      string  false  "Substring that can be contained in either the User's Display Name or Nickname"
// @Param        location   query      string  false  "User Location"
// @Param        is_verified   query      string  false  "User verification status"
// @Success      200  {object} 	string "Important Note: OK responses are wrapped in {"data": ... }"
// @Failure      400  {object} 	contracts.ErrResponse
// @Failure      404  {object} 	contracts.ErrResponse
// @Failure      500  {object}  contracts.ErrResponse
// @Router       /{version}/users/	[get]
func (h GetUsers) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req users2.GetUsersRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		if req.Nickname != "" {
			user, err := h.users.GetUserByNickname(ctx, req.Nickname)
			if err != nil {
				if errors.Is(err, contracts.ErrUserNotFound) {
					ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
					return
				}
				ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
				return
			}
			ctx.JSON(http.StatusOK, contracts.FormatOkResponse(user))
			return
		}

		req.Pagination.Validate()
		resUsers, err := h.users.GetUsers(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(resUsers))
	}
}
