package contracts

type UserSigninRequest struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UserSigninResponse struct {
	Id          string `json:"id" bson:"id"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Email       string `json:"email" bson:"email"`
	Type        string `json:"type" bson:"type"`
	AccessToken string `json:"access_token"`
}
