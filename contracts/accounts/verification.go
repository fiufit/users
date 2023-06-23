package accounts

type SendVerificationPinRequest struct {
	UserID      string
	PhoneNumber string `json:"phone_number" binding:"required"`
}
