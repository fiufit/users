package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const versionKey = "version"

type VersionHandlers map[string]gin.HandlerFunc

func HandleByVersion(handlers VersionHandlers) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if f, ok := handlers[ctx.Param(versionKey)]; ok {
			f(ctx)
		} else if f, ok := handlers[""]; ok {
			f(ctx)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{})
		}
	}
}
