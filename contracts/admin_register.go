package contracts

import "github.com/fiufit/users/models"

type AdminRegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type AdminLoginRequest AdminRegisterRequest

type AdminRegisterResponse struct {
	Admin models.Administrator `json:"admin"`
}

type AdminLoginResponse struct {
	Token string `json:"jwt"`
}
