package users

import (
	"github.com/fiufit/users/contracts/metrics"
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/repositories/external"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type UserFollower interface {
	FollowUser(ctx context.Context, req users.FollowUserRequest) error
	UnfollowUser(ctx context.Context, req users.UnfollowUserRequest) error
}

type UserFollowerImpl struct {
	notifications external.Notifications
	users         repositories.Users
	metrics       external.Metrics
	logger        *zap.Logger
}

func NewUserFollowerImpl(users repositories.Users, notifications external.Notifications, metrics external.Metrics, logger *zap.Logger) UserFollowerImpl {
	return UserFollowerImpl{users: users, notifications: notifications, metrics: metrics, logger: logger}
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

		followMetric := metrics.CreateMetricRequest{
			MetricType: "user_followed",
			SubType:    followedUser.ID,
		}
		uc.metrics.Create(ctx, followMetric)
	}
	return err
}

func (uc UserFollowerImpl) UnfollowUser(ctx context.Context, req users.UnfollowUserRequest) error {
	return uc.users.UnfollowUser(ctx, req.FollowedUserID, req.FollowerUserID)
}
