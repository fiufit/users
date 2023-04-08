package middleware

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/gin-gonic/gin"
)

type Nickname struct {
	Nickname string `form:"nickname" binding:"required"`
}

func BindNicknameFromQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var nick Nickname
		err := ctx.ShouldBindQuery(&nick)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}
		ctx.Set("nickname", nick.Nickname)
	}
}
