package users

import (
	"github.com/fiufit/users/contracts/users"
	"github.com/fiufit/users/repositories"
	"golang.org/x/net/context"
)

type UserFollower interface {
	FollowUser(ctx context.Context, req users.FollowUserRequest) error
	UnfollowUser(ctx context.Context, req users.UnfollowUserRequest) error
}

type UserFollowerImpl struct {
	users repositories.Users
}

func NewUserFollowerImpl(users repositories.Users) UserFollowerImpl {
	return UserFollowerImpl{users: users}
}

func (uc UserFollowerImpl) FollowUser(ctx context.Context, req users.FollowUserRequest) error {
	return uc.users.FollowUser(ctx, req.FollowedUserID, req.FollowerUserID)
}

func (uc UserFollowerImpl) UnfollowUser(ctx context.Context, req users.UnfollowUserRequest) error {
	return uc.users.UnfollowUser(ctx, req.FollowedUserID, req.FollowerUserID)
}
