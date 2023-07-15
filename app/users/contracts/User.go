package contracts

import "time"

type User struct {
	Id        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}
