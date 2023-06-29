package accounts

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	"github.com/fiufit/users/utils"
	mocks2 "github.com/fiufit/users/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
	"go.uber.org/zap/zaptest"
)

func TestSendVerificationPin_ErrUserAlreadyVerified(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.SendVerificationPinRequest{
		UserID:      "hola",
		PhoneNumber: "123",
	}

	auth.On("UserIsVerified", ctx, req.UserID).Return(true, nil)
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	_, err := verifier.SendVerificationPin(ctx, req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyVerified)
}

func TestSendVerificationPin_AuthErr(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.SendVerificationPinRequest{
		UserID:      "hola",
		PhoneNumber: "123",
	}

	auth.On("UserIsVerified", ctx, req.UserID).Return(false, errors.New("auth error"))
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	_, err := verifier.SendVerificationPin(ctx, req)
	assert.Error(t, err)
}

func TestSendVerificationPin_WhatsApperError(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.SendVerificationPinRequest{
		UserID:      "hola",
		PhoneNumber: "123",
	}

	baseTime := time.Now()

	timePatch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return baseTime
	})

	if err != nil {
		t.Fatalf("error patching time.Now: %v", err)
	}

	pin := models.VerificationPin{
		UserID:    req.UserID,
		Pin:       "$2a$10$TlyzU6dFSQBnv63YB.E6ieK1/ZkyX9IG9xT2zo1es6B/2YhwwJpaq", //hashed 1234
		ExpiresAt: baseTime.Add(time.Minute * 5),
	}

	patch, err := mpatch.PatchMethod(rand.Intn, func(num int) int {
		return 1234
	})

	if err != nil {
		t.Fatalf("error patching strconv.Itoa: %v", err)
	}

	hashPatch, err := mpatch.PatchMethod(utils.HashPassword, func(password string) (string, error) {
		return pin.Pin, nil
	})

	if err != nil {
		t.Fatalf("error patching utils.HashPassword: %v", err)
	}

	verification.On("Create", ctx, pin).Return(models.VerificationPin{}, nil)
	auth.On("UserIsVerified", ctx, req.UserID).Return(false, nil)
	whatsappSender.On("SendWhatsAppMessage", req.PhoneNumber, mock.Anything).Return(errors.New("whatsapp error"))
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	_, err = verifier.SendVerificationPin(ctx, req)
	patch.Unpatch()
	hashPatch.Unpatch()
	timePatch.Unpatch()
	assert.Error(t, err)
}

func TestSendVerificationPin_Ok(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.SendVerificationPinRequest{
		UserID:      "hola",
		PhoneNumber: "123",
	}

	baseTime := time.Now()

	timePatch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return baseTime
	})

	if err != nil {
		t.Fatalf("error patching time.Now: %v", err)
	}

	pin := models.VerificationPin{
		UserID:    req.UserID,
		Pin:       "$2a$10$TlyzU6dFSQBnv63YB.E6ieK1/ZkyX9IG9xT2zo1es6B/2YhwwJpaq", //hashed 1234
		ExpiresAt: baseTime.Add(time.Minute * 5),
	}

	patch, err := mpatch.PatchMethod(rand.Intn, func(num int) int {
		return 1234
	})

	if err != nil {
		t.Fatalf("error patching strconv.Itoa: %v", err)
	}

	hashPatch, err := mpatch.PatchMethod(utils.HashPassword, func(password string) (string, error) {
		return pin.Pin, nil
	})

	if err != nil {
		t.Fatalf("error patching utils.HashPassword: %v", err)
	}

	verification.On("Create", ctx, pin).Return(models.VerificationPin{}, nil)
	auth.On("UserIsVerified", ctx, req.UserID).Return(false, nil)
	whatsappSender.On("SendWhatsAppMessage", req.PhoneNumber, mock.Anything).Return(nil)
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	returnedPin, err := verifier.SendVerificationPin(ctx, req)
	patch.Unpatch()
	timePatch.Unpatch()
	hashPatch.Unpatch()
	assert.NoError(t, err)
	assert.Equal(t, returnedPin, pin)
}

func TestSendVerificationPin_VerificationCreateError(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.SendVerificationPinRequest{
		UserID:      "hola",
		PhoneNumber: "123",
	}

	baseTime := time.Now()

	timePatch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return baseTime
	})

	if err != nil {
		t.Fatalf("error patching time.Now: %v", err)
	}

	pin := models.VerificationPin{
		UserID:    req.UserID,
		Pin:       "$2a$10$TlyzU6dFSQBnv63YB.E6ieK1/ZkyX9IG9xT2zo1es6B/2YhwwJpaq", //hashed 1234
		ExpiresAt: baseTime.Add(time.Minute * 5),
	}

	patch, err := mpatch.PatchMethod(rand.Intn, func(num int) int {
		return 1234
	})

	if err != nil {
		t.Fatalf("error patching strconv.Itoa: %v", err)
	}

	_, err = mpatch.PatchMethod(utils.HashPassword, func(password string) (string, error) {
		return pin.Pin, nil
	})

	if err != nil {
		t.Fatalf("error patching utils.HashPassword: %v", err)
	}

	verification.On("Create", ctx, pin).Return(models.VerificationPin{}, errors.New("verification error"))
	auth.On("UserIsVerified", ctx, req.UserID).Return(false, nil)
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	_, err = verifier.SendVerificationPin(ctx, req)
	patch.Unpatch()
	timePatch.Unpatch()
	assert.Error(t, err)
}

func TestVerifyPin_RepoError(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.ValidateVerificationPinRequest{
		UserID: "hola",
		Pin:    "1234",
	}

	verification.On("GetByUserID", ctx, req.UserID).Return(models.VerificationPin{}, errors.New("repo error"))
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	err := verifier.VerifyPin(ctx, req)
	assert.Error(t, err)
}

func TestVerifyPin_PinDoesNotMatch(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.ValidateVerificationPinRequest{
		UserID: "hola",
		Pin:    "1234",
	}

	verification.On("GetByUserID", ctx, req.UserID).Return(models.VerificationPin{
		UserID: "hola",
		Pin:    "$2a$10$i85etNBEyrPPt9VQEOEAgub9RdBFsUKb.fpbZrDCE16Ti2OlKPVK.",
	}, nil)
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	err := verifier.VerifyPin(ctx, req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrInvalidVerificationPin)
}

func TestVerifyPin_PinExpired(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.ValidateVerificationPinRequest{
		UserID: "hola",
		Pin:    "1234",
	}

	verification.On("GetByUserID", ctx, req.UserID).Return(models.VerificationPin{
		UserID:    "hola",
		Pin:       "$2a$10$TlyzU6dFSQBnv63YB.E6ieK1/ZkyX9IG9xT2zo1es6B/2YhwwJpaq", //hashed 1234
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}, nil)
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	err := verifier.VerifyPin(ctx, req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, contracts.ErrVerificationPinExpired)
}

func TestVerifyPin_Ok(t *testing.T) {
	verification := new(mocks.VerificationPins)
	auth := new(mocks.Firebase)
	whatsappSender := new(mocks2.WhatsApper)
	logger := zaptest.NewLogger(t)

	ctx := context.Background()
	req := accounts.ValidateVerificationPinRequest{
		UserID: "hola",
		Pin:    "1234",
	}

	verification.On("GetByUserID", ctx, req.UserID).Return(models.VerificationPin{
		UserID:    "hola",
		Pin:       "$2a$10$TlyzU6dFSQBnv63YB.E6ieK1/ZkyX9IG9xT2zo1es6B/2YhwwJpaq", //hashed 1234
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}, nil)
	auth.On("VerifyUser", ctx, req.UserID).Return(nil)
	verifier := NewVerifierImpl(verification, auth, whatsappSender, logger)

	err := verifier.VerifyPin(ctx, req)
	assert.NoError(t, err)
}
