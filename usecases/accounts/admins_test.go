package accounts

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/accounts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories/mocks"
	utilMocks "github.com/fiufit/users/utils/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

func TestAdminLoginUnauthorizedError(t *testing.T) {
	adminRepo := new(mocks.Admins)
	toker := new(utilMocks.Toker)
	req := accounts.AdminLoginRequest{
		Email:    "testadmin@fiufit.com",
		Password: "wrongPassword",
	}
	ctx := context.Background()

	admin := models.Administrator{
		Model:    gorm.Model{},
		Email:    "testadmin@fiufit.com",
		Password: "$2a$10$gvDo.G4yR2T.Xdh.ZR9nouGnzXc4SjTbnFT3NBoJIFKxwBWoENXqa", //hunter2
	}

	adminRepo.On("GetByEmail", ctx, req.Email).Return(admin, nil)
	adminUc := NewAdminRegistererImpl(adminRepo, zaptest.NewLogger(t), toker)
	_, err := adminUc.Login(ctx, req)
	assert.ErrorIs(t, err, contracts.ErrInvalidPassword)
}

func TestAdminLoginNotFoundError(t *testing.T) {
	adminRepo := new(mocks.Admins)
	toker := new(utilMocks.Toker)
	req := accounts.AdminLoginRequest{
		Email:    "testadmin@fiufit.com",
		Password: "Password",
	}
	ctx := context.Background()

	adminRepo.On("GetByEmail", ctx, req.Email).Return(models.Administrator{}, contracts.ErrUserNotFound)
	adminUc := NewAdminRegistererImpl(adminRepo, zaptest.NewLogger(t), toker)
	_, err := adminUc.Login(ctx, req)
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestAdminLoginTokenCreationError(t *testing.T) {
	adminRepo := new(mocks.Admins)
	toker := new(utilMocks.Toker)
	req := accounts.AdminLoginRequest{
		Email:    "testadmin@fiufit.com",
		Password: "hunter2",
	}
	ctx := context.Background()
	admin := models.Administrator{
		Model:    gorm.Model{ID: 1},
		Email:    "testadmin@fiufit.com",
		Password: "$2a$10$gvDo.G4yR2T.Xdh.ZR9nouGnzXc4SjTbnFT3NBoJIFKxwBWoENXqa", //hunter2
	}

	adminRepo.On("GetByEmail", ctx, req.Email).Return(admin, nil)
	toker.On("CreateToken", strconv.Itoa(int(admin.ID)), true).Return("", errors.New("toker error"))
	adminUc := NewAdminRegistererImpl(adminRepo, zaptest.NewLogger(t), toker)

	_, err := adminUc.Login(ctx, req)
	assert.Error(t, err)
}

func TestAdminLoginOk(t *testing.T) {
	adminRepo := new(mocks.Admins)
	toker := new(utilMocks.Toker)
	req := accounts.AdminLoginRequest{
		Email:    "testadmin@fiufit.com",
		Password: "hunter2",
	}
	ctx := context.Background()
	admin := models.Administrator{
		Model:    gorm.Model{ID: 1},
		Email:    "testadmin@fiufit.com",
		Password: "$2a$10$gvDo.G4yR2T.Xdh.ZR9nouGnzXc4SjTbnFT3NBoJIFKxwBWoENXqa", //hunter2
	}

	adminRepo.On("GetByEmail", ctx, req.Email).Return(admin, nil)
	toker.On("CreateToken", strconv.Itoa(int(admin.ID)), true).Return("eyTokenCorrecto", nil)
	adminUc := NewAdminRegistererImpl(adminRepo, zaptest.NewLogger(t), toker)

	res, err := adminUc.Login(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, res.Token, "eyTokenCorrecto")
}

func TestAdminRegisterPasswordTooLongError(t *testing.T) {
	adminRepo := new(mocks.Admins)
	toker := new(utilMocks.Toker)
	req := accounts.AdminRegisterRequest{
		Email:    "testadmin@fiufit.com",
		Password: "passwordtoolongggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg",
	}
	ctx := context.Background()

	adminUc := NewAdminRegistererImpl(adminRepo, zaptest.NewLogger(t), toker)
	_, err := adminUc.Register(ctx, req)

	assert.Error(t, err)
}

/* TODO: find a way to make this test pass. As of now, if fails because the internal call to utils.HashPassword() is
non-deterministic (because we use bcrypt as hashing algorithm). In consequence, we can't mock adminRepo.Create() because
we don know the argument's hash field.
We could maybe convert utils.HashPassword() to a struct method, and mock an interface that that struct implements...

Until then, any tests involving Register() that reach the internal Create() call, will fail.
*/
func TestAdminRegisterRepoError(t *testing.T) {
	t.Skip("Read the TODO comment :(")
	adminRepo := new(mocks.Admins)
	toker := new(utilMocks.Toker)
	req := accounts.AdminRegisterRequest{
		Email:    "testadmin@fiufit.com",
		Password: "hunter2",
	}
	ctx := context.Background()
	admin := models.Administrator{
		Model:    gorm.Model{},
		Email:    req.Email,
		Password: "$2a$10$gvDo.G4yR2T.Xdh.ZR9nouGnzXc4SjTbnFT3NBoJIFKxwBWoENXqa",
	}

	adminUc := NewAdminRegistererImpl(adminRepo, zaptest.NewLogger(t), toker)
	adminRepo.On("Create", ctx, admin).Return(models.Administrator{}, errors.New("repo error"))
	_, err := adminUc.Register(ctx, req)

	assert.Error(t, err)
}
