package handlers

import (
	"errors"
	"net/http"

	"github.com/fiufit/users/contracts"
	acontracts "github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminLogin struct {
	admins accounts.AdminRegisterer
	logger *zap.Logger
}

func NewAdminLogin(admins accounts.AdminRegisterer, logger *zap.Logger) AdminLogin {
	return AdminLogin{admins: admins, logger: logger}
}

// Admin Login godoc
// @Summary      Lo gin as administrator
// @Description  Log in as administrator. Administrators and their credentials are created by other administrators
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        version   path      string  true  "API Version"
// @Param        payload   body acontracts.AdminLoginRequest true  "Body params"
// @Success      200  {object} 	accounts.AdminLoginResponse "Important Note: OK responses are wrapped in {"data": ... }"
// @Failure      400  {object} 	contracts.ErrResponse
// @Failure      401  {object} 	contracts.ErrResponse
// @Failure      404  {object}  contracts.ErrResponse
// @Failure      500  {object}  contracts.ErrResponse
// @Router       /{version}/admin/login 	[post]
func (h AdminLogin) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req acontracts.AdminLoginRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		res, err := h.admins.Login(ctx, req)
		if err != nil {
			if errors.Is(err, contracts.ErrInvalidPassword) {
				ctx.JSON(http.StatusUnauthorized, contracts.FormatErrResponse(contracts.ErrInvalidPassword))
				return

			} else if errors.Is(err, contracts.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, contracts.FormatErrResponse(contracts.ErrUserNotFound))
				return
			}

			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
