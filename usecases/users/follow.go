package users

import (
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/repositories"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type UserFollower interface {
	FollowUser(ctx context.Context, req users.FollowUserRequest) error
	UnfollowUser(ctx context.Context, req users.UnfollowUserRequest) error
}

type UserFollowerImpl struct {
	notifications repositories.Notifications
	users         repositories.Users
	logger        *zap.Logger
}

func NewUserFollowerImpl(users repositories.Users, notifications repositories.Notifications, logger *zap.Logger) UserFollowerImpl {
	return UserFollowerImpl{users: users, notifications: notifications, logger: logger}
}

func (uc UserFollowerImpl) FollowUser(ctx context.Context, req users.FollowUserRequest) error {
	followedUser, err := uc.users.GetByID(ctx, req.FollowedUserID)
	if err != nil {
		return err
	}

	followerUser, err := uc.users.GetByID(ctx, req.FollowerUserID)
	if err != nil {
		return err
	}
	err = uc.users.FollowUser(ctx, followedUser, followerUser)
	if err == nil {
		if uc.notifications.SendFollowersNotification(ctx, followerUser, followedUser) != nil {
			uc.logger.Error("Error sending notification", zap.Error(err))
		}
	}
	return err
}

func (uc UserFollowerImpl) UnfollowUser(ctx context.Context, req users.UnfollowUserRequest) error {
	return uc.users.UnfollowUser(ctx, req.FollowedUserID, req.FollowerUserID)
}
