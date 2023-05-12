package users

import (
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
)

type FollowUserRequest struct {
	FollowedUserID string
	FollowerUserID string `form:"follower_id" binding:"required"`
}

type UnfollowUserRequest struct {
	FollowedUserID string
	FollowerUserID string `uri:"followerID" binding:"required"`
}

type GetUserFollowersRequest struct {
	UserID string
	contracts.Pagination
}

type GetFollowedUsersRequest GetUserFollowersRequest

type GetUserFollowersResponse struct {
	contracts.Pagination
	Followers []models.User `json:"followers"`
}

type GetFollowedUsersResponse struct {
	contracts.Pagination
	Followed []models.User `json:"followed"`
}
