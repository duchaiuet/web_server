package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"web_server/middelwares"
	"web_server/model"
)

type controller struct {
	Service Service
}

func (c controller) Login(w http.ResponseWriter, r *http.Request) {
	var res LoginResponse
	var userPayload Login
	var user *model.User
	var tokenStr string

	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		res = LoginResponse{
			User:   nil,
			Mess:   err.Error(),
			Status: 400,
			Token:  "",
		}
		render.JSON(w, r, res)
		return
	}

	user, err = c.Service.GetUserByUserName(userPayload.UserName)
	if err != nil {
		res = LoginResponse{
			User:   nil,
			Mess:   err.Error(),
			Status: 500,
			Token:  "",
		}
		render.JSON(w, r, res)

		return
	}

	if user != nil {
		tokenStr, err = middelwares.GenerateToken(user)
	}
	if err != nil {
		res = LoginResponse{
			User:   nil,
			Mess:   err.Error(),
			Status: 500,
			Token:  "",
		}
		render.JSON(w, r, res)

		return
	}

	if middelwares.ComparePassword(userPayload.Password, user.Password) {
		res = LoginResponse{
			User:   user,
			Mess:   "oke",
			Status: 200,
			Token:  tokenStr,
		}
		render.JSON(w, r, res)

		return
	} else {
		res = LoginResponse{
			User:   nil,
			Mess:   "user not exist",
			Status: 500,
			Token:  "",
		}
		render.JSON(w, r, res)

		return
	}

}

func (c controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var res Response
	var user CreateUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res = Response{
			User:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	user.Password, err = middelwares.GeneratePassword(user.Password)
	if err != nil {
		res = Response{
			User:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}

	createUser, err := c.Service.CreateUser(user)
	if err != nil {
		res = Response{
			User:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}
	res = Response{
		User:   createUser,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) GetUser(w http.ResponseWriter, r *http.Request) {
	var res Response
	id := chi.URLParam(r, "id")
	if id == "" {
		res = Response{
			User:   nil,
			Mess:   "id is required",
			Status: 0,
		}

		render.JSON(w, r, res)
		return
	}

	user, err := c.Service.GetUserById(id)
	if err != nil {
		res = Response{
			User:   nil,
			Mess:   err.Error(),
			Status: 500,
		}

		render.JSON(w, r, res)
		return
	}
	res = Response{
		User:   user,
		Mess:   "ok",
		Status: 200,
	}
	render.JSON(w, r, res)
}

func (c controller) FilterUser(w http.ResponseWriter, r *http.Request) {
	var page, pageSize int
	var err error
	var res ListResponse

	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	sortBy := r.URL.Query().Get("sortBy")
	direction := r.URL.Query().Get("direction")
	status := r.URL.Query().Get("status")

	if page, err = strconv.Atoi(pageStr); err == nil {
		page = 1
	}

	if pageSize, err = strconv.Atoi(pageSizeStr); err == nil {
		pageSize = 10
	}

	searchFilter := GetFilterSearch(search, status, sortBy, direction, page, pageSize)
	users, err := c.Service.FilterUser(searchFilter)
	if err != nil {
		res = ListResponse{
			Mess:     err.Error(),
			Status:   500,
			Users:    nil,
			Page:     page,
			PageSize: pageSize,
		}

		render.JSON(w, r, res)
		return
	}

	res = ListResponse{
		Mess:     "ok",
		Status:   200,
		Users:    users,
		Page:     page,
		PageSize: pageSize,
	}

	render.JSON(w, r, res)
}

func (c controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var res Response
	var user UpdateUser

	id := chi.URLParam(r, "id")
	if id == "" {
		res = Response{
			User:   nil,
			Mess:   "id is required",
			Status: 0,
		}

		render.JSON(w, r, res)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res = Response{
			User:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	updateUser, err := c.Service.UpdateUser(user, id)
	if err != nil {
		res = Response{
			User:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)

		return
	}

	res = Response{
		User:   updateUser,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

type Controller interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	FilterUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewController(client *mongo.Client, collection string, database string) Controller {
	return controller{
		Service: NewService(client, collection, database),
	}
}
