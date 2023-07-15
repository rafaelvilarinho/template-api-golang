package libraries

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"api.template.com.br/helpers"
)

const (
	DATABASE = "template"
)

type SelectOptions struct {
	Sort  map[string]int
	Limit int64
	Skip  *int64
}

func SelectOne(collection string, query bson.M) (bson.M, error) {
	environment, _ := helpers.GetEnvironment()
	client, clientErr := mongo.Connect(context.TODO(), options.Client().ApplyURI(environment.MONGODB_URI))
	if clientErr != nil {
		return bson.M{}, nil
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(DATABASE).Collection(collection)

	var result bson.M
	err := coll.FindOne(context.TODO(), query).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return bson.M{}, err
	} else if err != nil {
		return bson.M{}, err
	}

	return result, nil
}

func SelectMany(collection string, query bson.M, opts ...SelectOptions) (*mongo.Cursor, error) {
	environment, _ := helpers.GetEnvironment()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(environment.MONGODB_URI))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(DATABASE).Collection(collection)

	findOpts := options.Find()

	if opts != nil {
		for _, opt := range opts {
			if len(opt.Sort) > 0 {
				sortRules := bson.D{}

				for key, val := range opt.Sort {
					sortRules = append(sortRules, bson.E{key, val})
				}

				findOpts.SetSort(sortRules)
			}

			if opt.Limit > 0 {
				findOpts.SetLimit(opt.Limit)
			}

			if opt.Skip != nil {
				findOpts.SetSkip(*opt.Skip)
			}
		}
	}

	if cursor, err := coll.Find(context.TODO(), query, findOpts); err != nil {
		return nil, err
	} else {
		return cursor, err
	}
}

func InsertOne(collection string, data bson.M) (*primitive.ObjectID, error) {
	environment, _ := helpers.GetEnvironment()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(environment.MONGODB_URI))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(DATABASE).Collection(collection)

	data["createdAt"] = time.Now()
	data["active"] = true

	if result, err := coll.InsertOne(context.TODO(), data); err != nil {
		return nil, err
	} else {
		if value, ok := result.InsertedID.(primitive.ObjectID); !ok {
			return nil, fmt.Errorf("error on get inserted id")
		} else {
			return &value, nil
		}
	}
}

func UpdateOne(collection string, query, data bson.M) (bool, error) {
	environment, _ := helpers.GetEnvironment()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(environment.MONGODB_URI))
	if err != nil {
		return false, err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(DATABASE).Collection(collection)

	if result, err := coll.UpdateOne(context.TODO(), query, data); err != nil {
		return false, err
	} else {
		return result.ModifiedCount > 0 && result.ModifiedCount == result.MatchedCount, nil
	}
}

func UpdateMany(collection string, query, data []bson.M) (bool, error) {
	environment, _ := helpers.GetEnvironment()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(environment.MONGODB_URI))
	if err != nil {
		return false, err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(DATABASE).Collection(collection)

	if result, err := coll.UpdateMany(context.TODO(), query, data); err != nil {
		return false, err
	} else {
		return result.ModifiedCount > 0 && result.ModifiedCount == result.MatchedCount, nil
	}
}
