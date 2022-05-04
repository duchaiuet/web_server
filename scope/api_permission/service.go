package api_permission

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"web_server/infrastructure"
	"web_server/model"
)

type service struct {
	Repository Repository
}

func (s service) CreateRule(rule CUCasbinRule) (*model.CasbinRule, error) {
	cRule := ConvertCUCasbinRuleToModel(rule)
	cRule.Id = primitive.NewObjectID()
	cRule.Type = "p"

	create, err := s.Repository.Create(cRule)
	if err != nil {
		return nil, err
	}

	return create, nil
}

func (s service) GetRuleById(id string) (*model.CasbinRule, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	permission, err := s.Repository.Get(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return permission, nil
}

func (s service) FilterRule(filter SearchFilter) ([]*model.CasbinRule, error) {
	permissions, err := s.Repository.Filter(filter)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return permissions, nil
}

func (s service) UpdateRule(rule CUCasbinRule, id string) (*model.CasbinRule, error) {
	update := ConvertCUCasbinRuleToModel(rule)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	update.Id = objId

	updatePermission, err := s.Repository.Update(update)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return updatePermission, nil
}

func (s service) DeleteRule(id string) error {
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
	CreateRule(rule CUCasbinRule) (*model.CasbinRule, error)
	GetRuleById(id string) (*model.CasbinRule, error)
	FilterRule(filter SearchFilter) ([]*model.CasbinRule, error)
	UpdateRule(rule CUCasbinRule, id string) (*model.CasbinRule, error)
	DeleteRule(id string) error
}

func NewService(client *mongo.Client, collection string, database string) Service {
	return service{
		Repository: NewRepository(client, collection, database),
	}
}
