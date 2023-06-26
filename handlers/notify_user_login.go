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

// Notify User Login godoc
//
//	@Summary		Creates a login metric for internal visualization.
//	@Description	Creates a login metric for internal visualization.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			version					path		string	true	"API Version"
//	@Param			method					query		string	true	"Login  method, either 'mail' or 'federated_entity'"
//	@Success		200						{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400						{object}	contracts.ErrResponse
//	@Router			/{version}/users/login	[post]
func (h NotifyUserLogin) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req accounts.NotifyLoginMethodRequest
		err := ctx.ShouldBindQuery(&req)
		validateErr := accounts.ValidateMethod(req.Method)
		if err != nil || validateErr != nil {
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
