package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"web_server/dto"
	"web_server/infrastructure"
	"web_server/model"
	"web_server/scope/role"
)

type service struct {
	UserRepository Repository
	RoleRepository role.Repository
}

func (s service) CreateUser(user dto.CreateUser) (*model.User, error) {
	createUser := dto.ConvertCreteUserToModel(user)

	tNow := time.Now()
	createUser.CreatedAt = &tNow
	createUser.Id = primitive.NewObjectID()

	create, err := s.UserRepository.Create(createUser)
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

	user, err := s.UserRepository.Get(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return user, nil
}

func (s service) GetUserByUserName(userName string) (*model.User, error) {

	user, err := s.UserRepository.GetByUser(userName)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	return user, nil
}

func (s service) FilterUser(filter dto.SearchFilter) ([]*model.User, int, error) {
	users, total, err := s.UserRepository.Filter(filter)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, 0, err
	}

	return users, total, nil
}

func (s service) UpdateUser(user dto.UpdateUser, id string) (*model.User, error) {
	update := dto.ConvertUpdateUserToModel(user)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return nil, err
	}

	tNow := time.Now()
	update.UpdatedAt = &tNow
	update.Id = objId

	updateUser, err := s.UserRepository.Update(update)
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

	err = s.UserRepository.Delete(objId)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return err
	}

	return nil
}

type Service interface {
	CreateUser(user dto.CreateUser) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	GetUserByUserName(userName string) (*model.User, error)
	FilterUser(filter dto.SearchFilter) ([]*model.User, int, error)
	UpdateUser(user dto.UpdateUser, id string) (*model.User, error)
	DeleteUser(id string) error
}

func NewService(client *mongo.Client, collection string, database string) Service {
	return service{
		UserRepository: NewUserRepository(client, collection, database),
	}
}
