package repositories

import (
	"context"
	"errors"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockery --name VerificationPins
type VerificationPins interface {
	Create(ctx context.Context, pin models.VerificationPin) (models.VerificationPin, error)
	GetByUserID(ctx context.Context, userID string) (models.VerificationPin, error)
}

type VerificationPinRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewVerificationPinRepository(db *gorm.DB, logger *zap.Logger) VerificationPinRepository {
	return VerificationPinRepository{db: db, logger: logger}
}

func (repo VerificationPinRepository) Create(ctx context.Context, pin models.VerificationPin) (models.VerificationPin, error) {
	db := repo.db.WithContext(ctx).Where("user_id = ?", pin.UserID)
	var newPin models.VerificationPin
	if err := db.Assign(&pin).FirstOrCreate(&newPin).Error; err != nil {
		repo.logger.Error("unable to create or update verification pin", zap.Error(err), zap.Any("pin", pin))
		return models.VerificationPin{}, err
	}
	return newPin, nil
}

func (repo VerificationPinRepository) GetByUserID(ctx context.Context, userID string) (models.VerificationPin, error) {
	db := repo.db.WithContext(ctx)

	var pin models.VerificationPin
	result := db.First(&pin, "user_id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.VerificationPin{}, contracts.ErrUserNotFound
		}
		return models.VerificationPin{}, result.Error
	}
	return pin, nil
}
