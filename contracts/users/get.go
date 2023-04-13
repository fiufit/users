package users

import "github.com/fiufit/users/contracts"

type GetUsersRequest struct {
	Name       string `form:"name"`
	Nickname   string `form:"nickname"`
	Location   string `form:"location"`
	IsVerified *bool  `form:"is_verified"`
	contracts.Pagination
}
