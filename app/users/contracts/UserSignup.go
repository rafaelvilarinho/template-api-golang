package contracts

type UserSignupRequest struct {
	FirstName      string `json:"firstName" bson:"firstName" binding:"required"`
	LastName       string `json:"lastName" bson:"lastName" binding:"required"`
	Email          string `json:"email" bson:"email" binding:"required"`
	Password       string `json:"password" bson:"password" binding:"required"`
	RetypePassword string `json:"retypePassword" bson:"retypePassword" binding:"required"`
}

type UserSignupResponse struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}
