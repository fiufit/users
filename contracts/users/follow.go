package users

import (
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
)

type FollowUserRequest struct {
	FollowedUserID string
	FollowerUserID string `form:"followerID" binding:"required"`
}

type UnfollowUserRequest struct {
	FollowedUserID string
	FollowerUserID string `uri:"followerID" binding:"required"`
}

type GetUserFollowersRequest struct {
	UserID string `form:"-"`
	contracts.Pagination
}

type GetUserFollowersResponse struct {
	contracts.Pagination
	Followers []models.User `json:"followers"`
}
