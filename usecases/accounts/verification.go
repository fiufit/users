package accounts

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/utils"
	"go.uber.org/zap"
)

type Verificator interface {
	SendVerificationPin(ctx context.Context, req accounts.SendVerificationPinRequest) (models.VerificationPin, error)
	VerifyPin(ctx context.Context, req accounts.ValidateVerificationPinRequest) error
}

type VerificatorImpl struct {
	verification   repositories.VerificationPins
	logger         *zap.Logger
	auth           repositories.Firebase
	whatsappSender utils.WhatsApper
}

func NewVerificatorImpl(verification repositories.VerificationPins, auth repositories.Firebase, whatsAppSender utils.WhatsApper, logger *zap.Logger) VerificatorImpl {
	return VerificatorImpl{logger: logger, auth: auth, verification: verification, whatsappSender: whatsAppSender}
}

func (uc *VerificatorImpl) SendVerificationPin(ctx context.Context, req accounts.SendVerificationPinRequest) (models.VerificationPin, error) {
	isVerified, err := uc.auth.UserIsVerified(ctx, req.UserID)
	if err != nil {
		return models.VerificationPin{}, err
	}
	if isVerified {
		return models.VerificationPin{}, contracts.ErrUserAlreadyVerified
	}
	rand.Seed(time.Now().UnixNano())
	randomPin := strconv.Itoa(rand.Intn(9999))
	hashedPin, err := utils.HashPassword(randomPin)
	if err != nil {
		return models.VerificationPin{}, err
	}
	pin := models.VerificationPin{
		UserID:    req.UserID,
		Pin:       hashedPin,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(5)),
	}
	_, err = uc.verification.Create(ctx, pin)
	if err != nil {
		return models.VerificationPin{}, err
	}
	err = uc.whatsappSender.SendWhatsAppMessage(req.PhoneNumber, randomPin)
	if err != nil {
		uc.logger.Error("unable to send whatsapp message", zap.Error(err), zap.Any("pin", pin))
		return models.VerificationPin{}, err
	}
	return pin, nil
}

func (uc *VerificatorImpl) VerifyPin(ctx context.Context, req accounts.ValidateVerificationPinRequest) error {
	pin, err := uc.verification.GetByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	if err := utils.ValidatePassword(req.Pin, pin.Pin); err != nil {
		return contracts.ErrInvalidVerificationPin
	}
	if time.Now().After(pin.ExpiresAt) {
		return contracts.ErrVerificationPinExpired
	}
	err = uc.auth.VerifyUser(ctx, req.UserID)
	return err
}
