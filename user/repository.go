package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"web_server/dto"
	"web_server/infrastructure"
	"web_server/model"
)

type Repository interface {
	Create(model.User) (*model.User, error)
	Get(id primitive.ObjectID) (*model.User, error)
	GetByUser(userName string) (*model.User, error)
	Filter(filter dto.SearchFilter) ([]*model.User, int, error)
	Update(model.User) (*model.User, error)
	Delete(id primitive.ObjectID) error
}

type repository struct {
	Db         *mongo.Database
	Collection string
}

func (r repository) GetByUser(userName string) (*model.User, error) {
	user := &model.User{}

	query := bson.D{{"user_name", userName}}
	err := r.Db.Collection(r.Collection).FindOne(context.TODO(), query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r repository) Create(user model.User) (*model.User, error) {
	_, err := r.Db.Collection(r.Collection).InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r repository) Get(id primitive.ObjectID) (*model.User, error) {
	user := &model.User{}

	query := bson.D{{"_id", id}}
	err := r.Db.Collection(r.Collection).FindOne(context.TODO(), query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r repository) Filter(filter dto.SearchFilter) ([]*model.User, int, error) {
	var users []*model.User
	query := bson.M{}

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.PageSize)).
		SetLimit(int64(filter.PageSize))

	if filter.SortBy.Field != "" {
		opts.SetSort(bson.D{{Key: filter.SortBy.Field, Value: -1}})
	}

	if filter.Search != "" {
		query = bson.M{
			"$or": bson.A{
				bson.M{"last_name": primitive.Regex{Pattern: ".*\\b" + filter.Search + "\\b.*", Options: "si"}},
				bson.M{"first_name": primitive.Regex{Pattern: ".*\\b" + filter.Search + "\\b.*", Options: "si"}},
			},
		}
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	cur, err := r.Db.Collection(r.Collection).Find(context.TODO(), query, opts)
	if err != nil {
		return nil, 0, err
	}
	total, err := r.Db.Collection(r.Collection).CountDocuments(context.TODO(), query)
	if err != nil {
		return nil, 0, err
	}

	infrastructure.InfoLog.Println("total: ", total)

	for cur.Next(context.TODO()) {
		var element model.User
		err := cur.Decode(&element)
		if err != nil {
			infrastructure.ErrLog.Println(err)
			break
		}

		users = append(users, &element)
	}

	return users, int(total), nil
}

func (r repository) Update(user model.User) (*model.User, error) {
	query := bson.D{{"_id", user.Id}}
	_, err := r.Db.Collection(r.Collection).UpdateOne(context.TODO(), query, bson.M{
		"$set": bson.M{
			"status":     user.Status,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.EmailAddress,
			"role":       user.Role,
		},
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r repository) Delete(id primitive.ObjectID) error {
	query := bson.D{{"_id", id}}
	_, err := r.Db.Collection(r.Collection).DeleteOne(context.TODO(), query)
	if err != nil {
		return err
	}

	return err
}

func NewUserRepository(client *mongo.Client, collection string, database string) Repository {
	return repository{
		Db:         client.Database(database),
		Collection: collection,
	}
}
