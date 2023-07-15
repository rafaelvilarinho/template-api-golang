package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryGeneralProps struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time         `json:"updatedAt" bson:"updatedAt"`
	DeletedAt *time.Time         `json:"deletedAt" bson:"deletedAt"`
}
