package accounts

type SendVerificationPinRequest struct {
	UserID      string
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type ValidateVerificationPinRequest struct {
	UserID string
	Pin    string `json:"pin" binding:"required"`
}
