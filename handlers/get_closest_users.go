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

type GetClosestUsers struct {
	users  users.UserGetter
	logger *zap.Logger
}

func NewGetClosestUsers(users users.UserGetter, logger *zap.Logger) GetClosestUsers {
	return GetClosestUsers{users: users, logger: logger}
}

func (h GetClosestUsers) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req users2.GetClosestUsersRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		req.UserID = ctx.MustGet("userID").(string)

		req.Pagination.Validate()
		resUsers, err := h.users.GetClosestUsers(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
				return
			}
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(resUsers))
	}
}
