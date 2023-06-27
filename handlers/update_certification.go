package handlers

import (
	"net/http"

	"github.com/fiufit/users/contracts"
	certContracts "github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/usecases/certifications"
	"github.com/gin-gonic/gin"
)

type UpdateCertification struct {
	certUpdater certifications.CertificationUpdater
}

func NewUpdateCertification(certUpdater certifications.CertificationUpdater) UpdateCertification {
	return UpdateCertification{certUpdater: certUpdater}
}

type certID struct {
	CertificationID uint `uri:"certificationID" binding:"required"`
}

func (h UpdateCertification) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req certContracts.UpdateCertificationRequest
		var cID certID
		err := ctx.ShouldBindUri(&cID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		err = ctx.ShouldBindQuery(&req)
		_, isCertStatusValid := models.ValidCertificationStatuses[req.Status]
		if err != nil || !isCertStatusValid {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}
		req.CertificationID = cID.CertificationID

		updatedCert, err := h.certUpdater.Update(ctx, req)
		if err != nil {
			contracts.HandleErrorType(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(updatedCert))
	}
}
