package contracts

type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required"`
}
