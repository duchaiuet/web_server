package api

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"web_server/model"
)

type Repository interface {
	Create(api model.Api) (*model.Api, error)
	Get(id primitive.ObjectID) (*model.Api, error)
	Filter(filter SearchFilter) ([]*model.Api, error)
	Update(api model.Api) (*model.Api, error)
	Delete(id primitive.ObjectID) error
}

type repository struct {
	Db         *mongo.Database
	Collection string
}

func (r repository) Create(api model.Api) (*model.Api, error) {
	_, err := r.Db.Collection(r.Collection).InsertOne(context.TODO(), api)
	if err != nil {
		return nil, err
	}
	return &api, err
}

func (r repository) Get(id primitive.ObjectID) (*model.Api, error) {
	api := &model.Api{}

	query := bson.D{{"_id", id}}
	err := r.Db.Collection(r.Collection).FindOne(context.TODO(), query).Decode(&api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func (r repository) Filter(filter SearchFilter) ([]*model.Api, error) {
	var users []*model.Api
	query := bson.M{}

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.PageSize)).
		SetLimit(int64(filter.PageSize))

	if filter.Search != "" {
		query = bson.M{
			"$or": bson.A{
				bson.M{"path": primitive.Regex{Pattern: filter.Search, Options: "si"}},
			},
		}
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	cur, err := r.Db.Collection(r.Collection).Find(context.TODO(), query, opts)
	if err != nil {
		return users, err
	}

	for cur.Next(context.TODO()) {
		var element model.Api
		err := cur.Decode(&element)
		if err != nil {
			return []*model.Api{}, err
		}

		users = append(users, &element)
	}

	return users, nil
}

func (r repository) Update(api model.Api) (*model.Api, error) {
	query := bson.D{{"_id", api.Id}}
	_, err := r.Db.Collection(r.Collection).UpdateOne(context.TODO(), query, bson.M{
		"$set": bson.M{
			"path":   api.Path,
			"status": api.Status,
		},
	})
	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (r repository) Delete(id primitive.ObjectID) error {
	query := bson.D{{"_id", id}}
	_, err := r.Db.Collection(r.Collection).DeleteOne(context.TODO(), query)
	if err != nil {
		return err
	}

	return err
}

func NewRoleRepository(client *mongo.Client, collection string, database string) Repository {
	return repository{
		Db:         client.Database(database),
		Collection: collection,
	}
}
