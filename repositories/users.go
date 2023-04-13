package repositories

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/database"
	"github.com/fiufit/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Users interface {
	GetByID(ctx context.Context, userID string) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
	Get(ctx context.Context, req users.GetUsersRequest) ([]models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
	auth   *auth.Client
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger, auth *auth.Client) UserRepository {
	return UserRepository{db: db, logger: logger, auth: auth}
}

func (repo UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	db := repo.db.WithContext(ctx)
	result := db.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.User{}, contracts.ErrUserAlreadyExists
		}
		repo.logger.Error("Unable to create user", zap.Error(result.Error), zap.Any("user", user))
		return models.User{}, result.Error
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
		repo.logger.Error("Unable to get user", zap.Error(result.Error), zap.String("ID", userID))
		return models.User{}, result.Error
	}

	return usr, nil
}

func (repo UserRepository) DeleteUser(ctx context.Context, userID string) error {
	db := repo.db.WithContext(ctx)
	var usr models.User
	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&usr, "id = ?", userID)
		if result.Error != nil {
			return result.Error
		}
		fbError := repo.auth.DeleteUser(ctx, userID)
		if fbError != nil {
			return fbError
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.ErrUserNotFound
		}
		repo.logger.Error("Unable to delete user", zap.Error(err), zap.String("ID", userID))
		return err
	}
	return nil
}
func (repo UserRepository) Get(ctx context.Context, req users.GetUsersRequest) ([]models.User, error) {
	var res []models.User
	db := repo.db.WithContext(ctx)
	if req.Name != "" {
		likeName := fmt.Sprintf("%%%v%%", req.Name)
		db = db.Where("LOWER(display_name) LIKE LOWER(?) OR LOWER(nickname) LIKE LOWER(?)", likeName, likeName)
	}
	if req.Location != "" {
		likeLocation := fmt.Sprintf("%%%v%%", req.Location)
		db = db.Where("LOWER(main_location) LIKE LOWER(?)", likeLocation)
	}
	if req.IsVerified != nil {
		db = db.Where("is_verified_trainer = ?", *req.IsVerified)
	}

	result := db.Scopes(database.Paginate(res, &req.Pagination, db)).Find(&res)
	return res, result.Error
}

func (repo UserRepository) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	db := repo.db.WithContext(ctx)
	var usr models.User
	result := db.Where("nickname = ?", nickname).First(&usr)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		repo.logger.Error("Unable to get user", zap.Error(result.Error), zap.String("nickname", nickname))
		return models.User{}, result.Error
	}

	return usr, nil
}

func (repo UserRepository) Update(ctx context.Context, user models.User) (models.User, error) {
	db := repo.db.WithContext(ctx)
	result := db.Save(&user)
	if result.Error != nil {
		repo.logger.Error("Unable to update user", zap.Error(result.Error), zap.Any("user", user))
		return models.User{}, result.Error
	}
	return user, nil
}
