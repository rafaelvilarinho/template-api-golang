package users

import (
	"context"
	"fmt"
	"time"

	"api.template.com.br/app/types"
	"api.template.com.br/libraries"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Logger *logrus.Logger
}

const (
	USER_COLLECTION = "users"
)

func (repository *UserRepository) FindAll(limit, skip int64, active *bool) ([]types.UserDTO, error) {
	var users []types.UserDTO

	query := bson.M{}

	if active != nil {
		query["active"] = active
	}

	opts := libraries.SelectOptions{}
	opts.Limit = limit
	opts.Skip = &skip

	if result, err := libraries.SelectMany(USER_COLLECTION, query, opts); err != nil {
		if err == mongo.ErrNoDocuments {
			return users, nil
		}

		return users, err
	} else {
		for result.Next(context.TODO()) {
			var userDTO types.UserDTO

			if bsonBytes, err := bson.Marshal(result.Current); err != nil {
				return nil, err
			} else {
				bson.Unmarshal(bsonBytes, &userDTO)

				users = append(users, userDTO)
			}
		}
	}

	return users, nil
}

func (repository *UserRepository) FindOneByEmail(email string) (*types.UserDTO, error) {
	if result, err := libraries.SelectOne(USER_COLLECTION, bson.M{
		"email": email,
	}); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	} else {
		var userDTO *types.UserDTO

		if _, filled := result["_id"].(primitive.ObjectID); !filled {
			return nil, nil
		} else {
			if bsonBytes, err := bson.Marshal(result); err != nil {
				return nil, err
			} else {
				bson.Unmarshal(bsonBytes, &userDTO)

				return userDTO, nil
			}
		}

	}
}

func (repository *UserRepository) Create(payload types.UserDTO) (*primitive.ObjectID, error) {
	user := bson.M{
		"firstName": payload.FirstName,
		"lastName":  payload.LastName,
		"email":     payload.Email,
		"password":  payload.Password,
		"type":      payload.Type,
		"verified":  payload.Verified,
	}

	if insertedId, err := libraries.InsertOne(USER_COLLECTION, user); err != nil {
		return nil, err
	} else {
		return insertedId, nil
	}
}

func (repository *UserRepository) VerifyAccount(userId string) (bool, error) {
	_userId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false, fmt.Errorf("user id parse error")
	}

	query := bson.M{
		"_id": _userId,
	}

	updateSet := bson.M{
		"$set": bson.M{
			"verified":  true,
			"updatedAt": time.Now(),
		},
	}

	if ok, err := libraries.UpdateOne(USER_COLLECTION, query, updateSet); err != nil {
		return false, err
	} else {
		return ok, nil
	}
}
