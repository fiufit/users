package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/contracts/metrics"
	"github.com/fiufit/users/repositories"
	"github.com/gin-gonic/gin"
)

type NotifyPasswordRecover struct {
	metrics repositories.Metrics
}

func NewNotifyPasswordRecover(metrics repositories.Metrics) NotifyPasswordRecover {
	return NotifyPasswordRecover{metrics: metrics}
}

func (h NotifyPasswordRecover) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req accounts.NotifyLoginMethodRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		metricsReq := metrics.CreateMetricRequest{
			MetricType: "password_recover",
		}

		h.metrics.Create(ctx, metricsReq)
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
