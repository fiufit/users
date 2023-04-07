package accounts

import (
	"context"
	"strconv"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminRegisterer interface {
	Login(ctx context.Context, req contracts.AdminLoginRequest) (contracts.AdminLoginResponse, error)
	Register(ctx context.Context, req contracts.AdminRegisterRequest) (contracts.AdminRegisterResponse, error)
}

type AdminRegistererImpl struct {
	admins repositories.Admins
	logger *zap.Logger
	toker  utils.JwtToker
}

func (uc *AdminRegistererImpl) Login(ctx context.Context, req contracts.AdminLoginRequest) (contracts.AdminLoginResponse, error) {
	admin, err := uc.admins.GetByEmail(ctx, req.Email)
	if err != nil {
		return contracts.AdminLoginResponse{}, err
	}

	if err := utils.ValidatePassword(req.Password, admin.Password); err != nil {
		return contracts.AdminLoginResponse{}, contracts.ErrInvalidPassword
	}

	token, err := uc.toker.CreateToken(strconv.Itoa(int(admin.ID)), true)
	if err != nil {
		uc.logger.Error("Unable to generate JWT for admin", zap.Error(err), zap.Any("admin", admin))
		return contracts.AdminLoginResponse{}, err
	}

	return contracts.AdminLoginResponse{Token: token}, nil
}

func (uc *AdminRegistererImpl) Register(ctx context.Context, req contracts.AdminRegisterRequest) (contracts.AdminRegisterResponse, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return contracts.AdminRegisterResponse{}, err
	}

	admin := models.Administrator{
		Model:    gorm.Model{},
		Email:    req.Email,
		Password: hashedPassword,
	}

	createdAdmin, err := uc.admins.Create(ctx, admin)
	if err != nil {
		return contracts.AdminRegisterResponse{}, err
	}

	return contracts.AdminRegisterResponse{Admin: createdAdmin}, nil
}
