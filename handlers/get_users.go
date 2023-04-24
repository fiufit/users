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
