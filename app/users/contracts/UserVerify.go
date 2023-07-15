package contracts

type UserVerifyRequest struct {
	VerifyCode string `json:"verifyCode" binding:"required"`
}
