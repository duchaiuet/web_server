package role

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

func (s service) GetRoleByCode(code string) (*model.Role, error) {
	role, err := s.Repository.GetByCode(code)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return role, nil
}

func (s service) CreateRole(role CURole) (*model.Role, error) {
	cRole := ConvertCURoleToModel(role)

	tNow := time.Now()
	cRole.CreatedAt = &tNow
	cRole.Id = primitive.NewObjectID()

	create, err := s.Repository.Create(cRole)
	if err != nil {
		return nil, err
	}

	return create, nil
}

func (s service) GetRoleById(id string) (*model.Role, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	role, err := s.Repository.Get(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return role, nil
}

func (s service) FilterRole(filter SearchFilter) ([]*model.Role, error) {
	roles, err := s.Repository.Filter(filter)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return roles, nil
}

func (s service) UpdateRole(role CURole, id string) (*model.Role, error) {
	update := ConvertCURoleToModel(role)

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

func (s service) DeleteRole(id string) error {
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
	CreateRole(role CURole) (*model.Role, error)
	GetRoleById(id string) (*model.Role, error)
	GetRoleByCode(code string) (*model.Role, error)
	FilterRole(filter SearchFilter) ([]*model.Role, error)
	UpdateRole(role CURole, id string) (*model.Role, error)
	DeleteRole(id string) error
}

func NewService(client *mongo.Client, collection string, database string) Service {
	return service{
		Repository: NewRoleRepository(client, collection, database),
	}
}
