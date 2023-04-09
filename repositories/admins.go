package repositories

import (
	"context"
	"errors"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Admins interface {
	GetByEmail(ctx context.Context, email string) (models.Administrator, error)
	Create(ctx context.Context, admin models.Administrator) (models.Administrator, error)
}

type AdminRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAdminRepository(db *gorm.DB, logger *zap.Logger) AdminRepository {
	return AdminRepository{db: db, logger: logger}
}

func (repo AdminRepository) Create(ctx context.Context, admin models.Administrator) (models.Administrator, error) {
	db := repo.db.WithContext(ctx)
	result := db.Create(&admin)
	if result.Error != nil {
		// TODO: check if this is error resulted from being a duplicate user in the database.
		return models.Administrator{}, errors.New("unable to create user")
	}
	return admin, nil
}

func (repo AdminRepository) GetByEmail(ctx context.Context, email string) (models.Administrator, error) {
	db := repo.db.WithContext(ctx)

	var admin models.Administrator
	result := db.First(&admin, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Administrator{}, contracts.ErrUserNotFound
		}
		return models.Administrator{}, result.Error
	}
	return admin, nil
}
