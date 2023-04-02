package middleware

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/gin-gonic/gin"
)

type UserID struct {
	UserID string `uri:"userID" binding:"required"`
}

func BindUserIDFromUri() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u UserID
		err := ctx.ShouldBindUri(&u)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}
		ctx.Set("userID", u.UserID)
	}
}
