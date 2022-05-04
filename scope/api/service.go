package api

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

func (s service) CreateAPI(api CUAPI) (*model.Api, error) {
	cAPI := ConvertCUAPIToModel(api)

	tNow := time.Now()
	cAPI.CreatedAt = &tNow
	cAPI.Id = primitive.NewObjectID()

	create, err := s.Repository.Create(cAPI)
	if err != nil {
		return nil, err
	}

	return create, nil
}

func (s service) GetAPIById(id string) (*model.Api, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	api, err := s.Repository.Get(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return api, nil
}

func (s service) FilterAPI(filter SearchFilter) ([]*model.Api, error) {
	apis, err := s.Repository.Filter(filter)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return apis, nil
}

func (s service) UpdateAPI(api CUAPI, id string) (*model.Api, error) {
	update := ConvertCUAPIToModel(api)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	tNow := time.Now()
	update.UpdatedAt = &tNow
	update.Id = objId

	updateApi, err := s.Repository.Update(update)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return updateApi, nil
}

func (s service) DeleteAPI(id string) error {
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
	CreateAPI(api CUAPI) (*model.Api, error)
	GetAPIById(id string) (*model.Api, error)
	FilterAPI(filter SearchFilter) ([]*model.Api, error)
	UpdateAPI(role CUAPI, id string) (*model.Api, error)
	DeleteAPI(id string) error
}

func NewService(client *mongo.Client, collection string, database string) Service {
	return service{
		Repository: NewRoleRepository(client, collection, database),
	}
}
