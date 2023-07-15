package types

type UserDTO struct {
	RepositoryGeneralProps `json:",inline" bson:",inline"`
	FirstName              string `json:"firstName" bson:"firstName"`
	LastName               string `json:"lastName" bson:"lastName"`
	Email                  string `json:"email" bson:"email"`
	Password               string `json:"password" bson:"password"`
	Avatar                 string `json:"avatar" bson:"avatar"`
	Gender                 string `json:"gender" bson:"gender"`
	Type                   string `json:"type" bson:"type"`
	Verified               bool   `json:"verified" bson:"verified"`
}
