package users

import (
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
)

type GetUsersRequest struct {
	Name       string `form:"name"`
	Nickname   string `form:"nickname"`
	IsVerified *bool  `form:"is_verified"`
	Disabled   *bool  `form:"disabled"`
	contracts.Pagination
}

type GetClosestUsersRequest struct {
	UserID    string
	Latitude  float64
	Longitude float64
	Distance  uint `form:"distance" binding:"required"`
	contracts.Pagination
}

type GetUsersResponse struct {
	Pagination contracts.Pagination `json:"pagination"`
	Users      []models.User        `json:"users"`
}
