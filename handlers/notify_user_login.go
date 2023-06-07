package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/contracts/metrics"
	"github.com/fiufit/users/repositories"
	"github.com/gin-gonic/gin"
)

type NotifyUserLogin struct {
	metrics repositories.Metrics
}

func NewNotifyUserLogin(metrics repositories.Metrics) NotifyUserLogin {
	return NotifyUserLogin{metrics: metrics}
}

func (h NotifyUserLogin) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req accounts.NotifyLoginMethodRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		metricsReq := metrics.CreateMetricRequest{
			MetricType: "login",
			SubType:    req.Method,
		}

		h.metrics.Create(ctx, metricsReq)
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
