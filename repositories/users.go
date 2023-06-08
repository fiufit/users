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
	"github.com/fiufit/users/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockery --name Users
type Users interface {
	GetByID(ctx context.Context, userID string) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
	Get(ctx context.Context, req ucontracts.GetUsersRequest) (ucontracts.GetUsersResponse, error)
	GetByDistance(ctx context.Context, req ucontracts.GetClosestUsersRequest) (ucontracts.GetUsersResponse, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, userID string) error
	FollowUser(ctx context.Context, followedUserID string, followerUserID string) error
	UnfollowUser(ctx context.Context, followedUserID string, followerUserID string) error
	GetFollowers(ctx context.Context, request ucontracts.GetUserFollowersRequest) (ucontracts.GetUserFollowersResponse, error)
	GetFollowed(ctx context.Context, req ucontracts.GetFollowedUsersRequest) (ucontracts.GetFollowedUsersResponse, error)
}

type UserRepository struct {
	db             *gorm.DB
	logger         *zap.Logger
	auth           Firebase
	reverseLocator *utils.ReverseLocator
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger, auth Firebase, reverseLocator *utils.ReverseLocator) UserRepository {
	return UserRepository{db: db, logger: logger, auth: auth, reverseLocator: reverseLocator}
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
	repo.fillUserLocation(&user)
	return user, nil
}

func (repo UserRepository) GetByID(ctx context.Context, userID string) (models.User, error) {
	db := repo.db.WithContext(ctx)
	var usr models.User
	result := db.Preload("Interests").First(&usr, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		repo.logger.Error("Unable to get user", zap.Error(result.Error), zap.String("ID", userID))
		return models.User{}, result.Error
	}

	repo.fillUserLocation(&usr)
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

	if req.IsVerified != nil {
		db = db.Where("is_verified_trainer = ?", *req.IsVerified)
	}
	if req.Disabled != nil {
		db = db.Where("disabled = ?", *req.Disabled)
	}
	if req.Name != "" {
		likeName := fmt.Sprintf("%v%%", strings.ToLower(req.Name))
		db = db.Where("LOWER(display_name) LIKE ? OR LOWER(nickname) LIKE ?", likeName, likeName)
	}

	result := db.Scopes(database.Paginate(res, &req.Pagination, db)).Preload("Interests").Find(&res)
	if result.Error != nil {
		repo.logger.Error("Unable to get users with pagination", zap.Error(result.Error), zap.Any("request", req))
		return ucontracts.GetUsersResponse{}, result.Error
	}
	for i := range res {
		repo.fillUserLocation(&res[i])
	}

	return ucontracts.GetUsersResponse{Users: res, Pagination: req.Pagination}, nil
}

func (repo UserRepository) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	db := repo.db.WithContext(ctx)
	var usr models.User
	result := db.Where("nickname = ?", nickname).Preload("Interests").First(&usr)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, contracts.ErrUserNotFound
		}
		repo.logger.Error("Unable to get user", zap.Error(result.Error), zap.String("nickname", nickname))
		return models.User{}, result.Error
	}

	repo.fillUserLocation(&usr)
	return usr, nil
}

func (repo UserRepository) GetByDistance(ctx context.Context, req ucontracts.GetClosestUsersRequest) (ucontracts.GetUsersResponse, error) {
	db := repo.db.WithContext(ctx)
	var closestUsers []models.User

	// TODO: Find out how to order by earthdistance too using gorm
	result := db.
		Scopes(database.Paginate(closestUsers, &req.Pagination, db)).
		Where("earth_distance(ll_to_earth(?, ?), ll_to_earth(users.latitude, users.longitude)) <= ? AND users.ID != ?", req.Latitude, req.Longitude, req.Distance*1000, req.UserID).
		Preload("Interests").
		Find(&closestUsers)

	if result.Error != nil {
		repo.logger.Error("Unable to get closest users with pagination", zap.Error(result.Error), zap.Any("request", req))
		return ucontracts.GetUsersResponse{}, result.Error
	}

	for i := range closestUsers {
		repo.fillUserLocation(&closestUsers[i])
	}

	return ucontracts.GetUsersResponse{Users: closestUsers, Pagination: req.Pagination}, nil
}

func (repo UserRepository) Update(ctx context.Context, user models.User) (models.User, error) {
	db := repo.db.WithContext(ctx)

	err := db.Model(&user).Association("Interests").Replace(user.Interests)
	if err != nil {
		repo.logger.Error("Unable to update user interests", zap.Error(err), zap.Any("user", user))
		return models.User{}, err
	}

	result := db.Save(&user)
	if result.Error != nil {
		repo.logger.Error("Unable to update user", zap.Error(result.Error), zap.Any("user", user))
		return models.User{}, result.Error
	}
	repo.fillUserLocation(&user)
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

	// TODO: figure out how to do this properly inside database.Paginate(). We have to overwrite the totalRows with
	// the following count(), because otherwise the total user count is set. Maybe use db.Scopes() ?

	req.Pagination.TotalRows = db.Model(&user).Association("Followers").Count()

	for i := range user.Followers {
		repo.fillUserLocation(&user.Followers[i])
	}

	response := ucontracts.GetUserFollowersResponse{
		Pagination: req.Pagination,
		Followers:  user.Followers,
	}

	return response, nil
}

func (repo UserRepository) GetFollowed(ctx context.Context, req ucontracts.GetFollowedUsersRequest) (ucontracts.GetFollowedUsersResponse, error) {
	db := repo.db.WithContext(ctx)
	var followedUsers []models.User

	db = db.Model(&followedUsers).Joins("LEFT JOIN user_followers ON user_followers.user_id = users.id").Where("user_followers.follower_id = ?", req.UserID)
	result := db.Scopes(database.Paginate(followedUsers, &req.Pagination, db)).Find(&followedUsers)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ucontracts.GetFollowedUsersResponse{}, contracts.ErrUserNotFound
		}

		repo.logger.Error("unable to get user followers", zap.Error(result.Error), zap.String("userID", req.UserID))
		return ucontracts.GetFollowedUsersResponse{}, result.Error
	}

	for i := range followedUsers {
		repo.fillUserLocation(&followedUsers[i])
	}

	response := ucontracts.GetFollowedUsersResponse{
		Pagination: req.Pagination,
		Followed:   followedUsers,
	}

	return response, nil
}

func (repo UserRepository) fillUserLocation(user *models.User) {
	usrLocation, err := repo.reverseLocator.GetLocationFromCoordinates(user.Latitude, user.Longitude)
	if err != nil {
		repo.logger.Error("Unable to reverse geolocate user's coordinates")
		return
	}
	user.MainLocation = usrLocation
}
