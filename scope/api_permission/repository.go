package api_permission

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"web_server/infrastructure"
	"web_server/model"
)

type Repository interface {
	Create(permission model.CasbinRule) (*model.CasbinRule, error)
	Get(id primitive.ObjectID) (*model.CasbinRule, error)
	Filter(filter SearchFilter) ([]*model.CasbinRule, error)
	Update(permission model.CasbinRule) (*model.CasbinRule, error)
	Delete(id primitive.ObjectID) error
}

type repository struct {
	Db         *mongo.Database
	Collection string
}

func (r repository) Create(rule model.CasbinRule) (*model.CasbinRule, error) {
	_, err := r.Db.Collection(r.Collection).InsertOne(context.TODO(), rule)
	if err != nil {
		return nil, err
	}
	return &rule, err
}

func (r repository) Get(id primitive.ObjectID) (*model.CasbinRule, error) {
	rule := &model.CasbinRule{}

	query := bson.D{{"_id", id}}
	err := r.Db.Collection(r.Collection).FindOne(context.TODO(), query).Decode(&rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (r repository) Filter(filter SearchFilter) ([]*model.CasbinRule, error) {
	var permissions []*model.CasbinRule
	query := bson.M{}

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.PageSize)).
		SetLimit(int64(filter.PageSize))

	if filter.Search != "" {
		query = bson.M{
			"$or": bson.A{
				bson.M{"v0": primitive.Regex{Pattern: filter.Search, Options: "si"}},
				bson.M{"v1": primitive.Regex{Pattern: filter.Search, Options: "si"}},
			},
		}
	}

	cur, err := r.Db.Collection(r.Collection).Find(context.TODO(), query, opts)
	if err != nil {
		return permissions, err
	}

	for cur.Next(context.TODO()) {
		var element model.CasbinRule
		err = cur.Decode(&element)
		if err != nil {
			infrastructure.ErrLog.Println(err)
			break
		}

		permissions = append(permissions, &element)
	}

	return permissions, nil
}

func (r repository) Update(rule model.CasbinRule) (*model.CasbinRule, error) {
	query := bson.D{{"_id", rule.Id}}
	_, err := r.Db.Collection(r.Collection).UpdateOne(context.TODO(), query, bson.M{
		"$set": bson.M{
			"v1": rule.Path,
			"v2": rule.Rule,
		},
	})
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (r repository) Delete(id primitive.ObjectID) error {
	query := bson.D{{"_id", id}}
	_, err := r.Db.Collection(r.Collection).DeleteOne(context.TODO(), query)
	if err != nil {
		return err
	}

	return err
}

func NewRepository(client *mongo.Client, collection string, database string) Repository {
	return repository{
		Db:         client.Database(database),
		Collection: collection,
	}
}
