package contracts

type AdminRegisterRequest struct {
	Email    string `json:"email" binding:"required;email"`
	Password string `json:"password" binding:"required"`
}

type AdminRegisterResponse struct {
	Token string `json:"jwt"`
}

type AdminPasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
