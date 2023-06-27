package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	certContracts "github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/usecases/certifications"
	"github.com/gin-gonic/gin"
)

type CreateCertification struct {
	certCreator certifications.CertificationCreator
}

func NewCreateCertification(certCreator certifications.CertificationCreator) CreateCertification {
	return CreateCertification{certCreator: certCreator}
}

func (h CreateCertification) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req certContracts.CreateCertificationRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		cert, err := h.certCreator.Create(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(cert))
	}
}
