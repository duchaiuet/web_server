package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"web_server/infrastructure"
	"web_server/model"
)

type service struct {
	Repository Repository
}

func (s service) CreateUser(user CreateUser) (*model.User, error) {
	createUser := ConvertCreteUserToModel(user)

	tNow := time.Now()
	createUser.CreatedAt = &tNow
	createUser.Id = primitive.NewObjectID()

	create, err := s.Repository.Create(createUser)
	if err != nil {
		return nil, err
	}

	return create, nil
}

func (s service) GetUserById(id string) (*model.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	user, err := s.Repository.Get(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return user, nil
}

func (s service) GetUserByUserName(userName string) (*model.User, error) {

	user, err := s.Repository.GetByUser(userName)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return user, nil
}

func (s service) FilterUser(filter SearchFilter) ([]*model.User, error) {
	users, err := s.Repository.Filter(filter)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return users, nil
}

func (s service) UpdateUser(user UpdateUser, id string) (*model.User, error) {
	update := ConvertUpdateUserToModel(user)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	tNow := time.Now()
	update.UpdatedAt = &tNow
	update.Id = objId

	updateUser, err := s.Repository.Update(update)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return updateUser, nil
}

func (s service) DeleteUser(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return err
	}

	err = s.Repository.Delete(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return err
	}

	return nil
}

type Service interface {
	CreateUser(user CreateUser) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	GetUserByUserName(userName string) (*model.User, error)
	FilterUser(filter SearchFilter) ([]*model.User, error)
	UpdateUser(user UpdateUser, id string) (*model.User, error)
	DeleteUser(id string) error
}

func NewService(client *mongo.Client, collection string, database string) Service {
	return service{
		Repository: NewUserRepository(client, collection, database),
	}
}
