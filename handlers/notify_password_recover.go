package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
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

// Notify Password Recovery godoc
//
//	@Summary		Creates a password recovery metric for internal visualization.
//	@Description	Creates a password recovery metric for internal visualization.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string	true	"API Version"
//	@Success		200									{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Router			/{version}/users/password-recover 	[post]
func (h NotifyPasswordRecover) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		metricsReq := metrics.CreateMetricRequest{
			MetricType: "password_recover",
		}

		h.metrics.Create(ctx, metricsReq)
		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
