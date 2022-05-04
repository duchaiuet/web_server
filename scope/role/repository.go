package role

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"web_server/model"
)

type Repository interface {
	Create(role model.Role) (*model.Role, error)
	Get(id primitive.ObjectID) (*model.Role, error)
	GetByCode(code string) (*model.Role, error)
	Filter(filter SearchFilter) ([]*model.Role, error)
	Update(role model.Role) (*model.Role, error)
	Delete(id primitive.ObjectID) error
}

type repository struct {
	Db         *mongo.Database
	Collection string
}

func (r repository) GetByCode(code string) (*model.Role, error) {
	role := &model.Role{}

	query := bson.D{{"code", code}}
	err := r.Db.Collection(r.Collection).FindOne(context.TODO(), query).Decode(&role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r repository) Create(role model.Role) (*model.Role, error) {
	_, err := r.Db.Collection(r.Collection).InsertOne(context.TODO(), role)
	if err != nil {
		return nil, err
	}
	return &role, err
}

func (r repository) Get(id primitive.ObjectID) (*model.Role, error) {
	role := &model.Role{}

	query := bson.D{{"_id", id}}
	err := r.Db.Collection(r.Collection).FindOne(context.TODO(), query).Decode(&role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r repository) Filter(filter SearchFilter) ([]*model.Role, error) {
	var users []*model.Role
	query := bson.M{}

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.PageSize)).
		SetLimit(int64(filter.PageSize))

	if filter.Search != "" {
		query = bson.M{
			"$or": bson.A{
				bson.M{"name": primitive.Regex{Pattern: ".*\\b" + filter.Search + "\\b.*", Options: "si"}},
				bson.M{"code": primitive.Regex{Pattern: filter.Search, Options: "si"}},
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
		var element model.Role
		err := cur.Decode(&element)
		if err != nil {
			return []*model.Role{}, err
		}

		users = append(users, &element)
	}

	return users, nil
}

func (r repository) Update(role model.Role) (*model.Role, error) {
	query := bson.D{{"_id", role.Id}}
	_, err := r.Db.Collection(r.Collection).UpdateOne(context.TODO(), query, bson.M{
		"$set": bson.M{
			"code":   role.Code,
			"name":   role.Name,
			"status": role.Status,
		},
	})
	if err != nil {
		return nil, err
	}

	return &role, nil
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
