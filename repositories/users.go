package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fiufit/users/contracts"
	ucontracts "github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/database"
	"github.com/fiufit/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockery --name Users
type Users interface {
	GetByID(ctx context.Context, userID string) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
	Get(ctx context.Context, req ucontracts.GetUsersRequest) (ucontracts.GetUsersResponse, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, userID string) error
	FollowUser(ctx context.Context, followedUserID string, followerUserID string) error
	UnfollowUser(ctx context.Context, followedUserID string, followerUserID string) error
	GetFollowers(ctx context.Context, request ucontracts.GetUserFollowersRequest) (ucontracts.GetUserFollowersResponse, error)
}

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
	auth   Firebase
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger, auth Firebase) UserRepository {
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
func (repo UserRepository) Get(ctx context.Context, req ucontracts.GetUsersRequest) (ucontracts.GetUsersResponse, error) {
	var res []models.User
	db := repo.db.WithContext(ctx)

	if req.Location != "" {
		likeLocation := fmt.Sprintf("%%%v%%", req.Location)
		db = db.Where("LOWER(main_location) LIKE LOWER(?)", likeLocation)
	}
	if req.IsVerified != nil {
		db = db.Where("is_verified_trainer = ?", *req.IsVerified)
	}

	if req.Name != "" {
		likeName := fmt.Sprintf("%v%%", strings.ToLower(req.Name))
		db = db.Where("LOWER(display_name) LIKE ? OR LOWER(nickname) LIKE ?", likeName, likeName)
	}

	result := db.Scopes(database.Paginate(res, &req.Pagination, db)).Find(&res)
	if result.Error != nil {
		repo.logger.Error("Unable to get users with pagination", zap.Error(result.Error), zap.Any("request", req))
		return ucontracts.GetUsersResponse{}, result.Error
	}

	return ucontracts.GetUsersResponse{Users: res, Pagination: req.Pagination}, nil
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

func (repo UserRepository) FollowUser(ctx context.Context, followedUserID string, followerUserID string) error {
	followedUser, err := repo.GetByID(ctx, followedUserID)
	if err != nil {
		return err
	}

	followerUser, err := repo.GetByID(ctx, followerUserID)
	if err != nil {
		return err
	}
	db := repo.db.WithContext(ctx)

	return db.Model(&followedUser).Association("Followers").Append(&followerUser)
}

func (repo UserRepository) UnfollowUser(ctx context.Context, followedUserID string, followerUserID string) error {
	followedUser, err := repo.GetByID(ctx, followedUserID)
	if err != nil {
		return err
	}

	followerUser, err := repo.GetByID(ctx, followerUserID)
	if err != nil {
		return err
	}
	db := repo.db.WithContext(ctx)

	return db.Model(&followedUser).Association("Followers").Delete(&followerUser)
}

func (repo UserRepository) GetFollowers(ctx context.Context, req ucontracts.GetUserFollowersRequest) (ucontracts.GetUserFollowersResponse, error) {
	db := repo.db.WithContext(ctx)
	var user models.User
	result := db.Preload("Followers", database.Paginate(&user.Followers, &req.Pagination, db)).First(&user, "users.id = ?", req.UserID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ucontracts.GetUserFollowersResponse{}, contracts.ErrUserNotFound
		}

		repo.logger.Error("unable to get user followers", zap.Error(result.Error), zap.String("userID", req.UserID))
		return ucontracts.GetUserFollowersResponse{}, result.Error
	}

	// TODO: figure out how to do this properly inside database.Paginate(). We have to overwrite the totalrows with
	// the following count(), because otherwise the total user count is set.

	req.Pagination.TotalRows = db.Model(&user).Association("Followers").Count()

	response := ucontracts.GetUserFollowersResponse{
		Pagination: req.Pagination,
		Followers:  user.Followers,
	}

	return response, nil
}
