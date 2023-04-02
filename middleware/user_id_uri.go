package middleware

import (
	"net/http"

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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.Set("userID", u.UserID)
	}
}
