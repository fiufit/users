package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	certContracts "github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/usecases/certifications"
	"github.com/gin-gonic/gin"
)

type GetCertifications struct {
	certGetter certifications.CertificationGetter
}

func NewGetCertifications(certGetter certifications.CertificationGetter) GetCertifications {
	return GetCertifications{certGetter: certGetter}
}

func (h GetCertifications) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req certContracts.GetCertificationsRequest
		err := ctx.ShouldBindQuery(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		res, err := h.certGetter.Get(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
