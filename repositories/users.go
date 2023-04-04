package repositories

import (
	"context"
	"errors"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Users interface {
	GetByID(ctx context.Context, userID string) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return UserRepository{db: db, logger: logger}
}

func (repo UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	db := repo.db.WithContext(ctx)
	result := db.Create(&user)
	if result.Error != nil {
		// TODO: check if this is error resulted from being a duplicate user in the database.
		return models.User{}, errors.New("unable to create user")
	}
	return user, nil
}

func (repo UserRepository) GetByID(ctx context.Context, userID string) (models.User, error) {
	db := repo.db.WithContext(ctx)
	var usr models.User
	result := db.First(&usr, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		return models.User{}, result.Error
	}

	return usr, nil
}

func (repo UserRepository) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	db := repo.db.WithContext(ctx)
	var usr models.User
	result := db.Where("nickname = ?", nickname).First(&usr)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		return models.User{}, result.Error
	}

	return usr, nil
}
