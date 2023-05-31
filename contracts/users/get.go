package users

import (
	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
)

type GetUsersRequest struct {
	Name       string `form:"name"`
	Nickname   string `form:"nickname"`
	Location   string `form:"location"`
	IsVerified *bool  `form:"is_verified"`
	Disabled   *bool  `form:"disabled"`
	contracts.Pagination
}

type GetUsersResponse struct {
	Pagination contracts.Pagination `json:"pagination"`
	Users      []models.User        `json:"users"`
}
