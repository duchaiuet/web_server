package user

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"web_server/dto"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"web_server/infrastructure"
	"web_server/middlewares"
	"web_server/model"
	"web_server/scope/role"
)

type controller struct {
	UserService Service
	RoleService role.Service
}

func (c controller) Register(w http.ResponseWriter, r *http.Request) {
	var res dto.Response
	var user dto.CreateUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	user.Password, err = middlewares.GeneratePassword(user.Password)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}

	createUser, err := c.UserService.CreateUser(user)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}
	userRes := dto.ConvertUserToResponse(createUser)
	res = dto.Response{
		Data:   &userRes,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) Login(w http.ResponseWriter, r *http.Request) {
	var res dto.LoginResponse
	var userPayload dto.Login
	var user *model.User
	var tokenStr string

	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		res = dto.LoginResponse{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
			Token:  "",
		}
		render.JSON(w, r, res)
		return
	}

	user, err = c.UserService.GetUserByUserName(userPayload.UserName)
	if err != nil {
		res = dto.LoginResponse{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
			Token:  "",
		}
		render.JSON(w, r, res)

		return
	}

	if user != nil {
		userRes := dto.ConvertUserToResponse(user)

		tokenStr, err = middlewares.GenerateToken(userRes)

		if err != nil {
			res = dto.LoginResponse{
				Data:   nil,
				Mess:   err.Error(),
				Status: 500,
				Token:  "",
			}
			render.JSON(w, r, res)

			return
		}

		if middlewares.ComparePassword(userPayload.Password, user.Password) {
			res = dto.LoginResponse{
				Data:   &userRes,
				Mess:   "oke",
				Status: 200,
				Token:  tokenStr,
			}
			render.JSON(w, r, res)

			return
		} else {
			res = dto.LoginResponse{
				Data:   nil,
				Mess:   "user not exist",
				Status: 500,
				Token:  "",
			}
			render.JSON(w, r, res)

			return
		}
	} else {
		res = dto.LoginResponse{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
			Token:  "",
		}
		render.JSON(w, r, res)

		return
	}

}

func (c controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var res dto.Response
	var user dto.CreateUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	user.Password, err = middlewares.GeneratePassword(user.Password)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}

	createUser, err := c.UserService.CreateUser(user)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}
	userRes := dto.ConvertUserToResponse(createUser)
	res = dto.Response{
		Data:   &userRes,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) GetUser(w http.ResponseWriter, r *http.Request) {
	var res dto.Response
	id := chi.URLParam(r, "id")
	if id == "" {
		res = dto.Response{
			Data:   nil,
			Mess:   "id is required",
			Status: 0,
		}

		render.JSON(w, r, res)
		return
	}

	user, err := c.UserService.GetUserById(id)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}

		render.JSON(w, r, res)
		return
	}

	userRes := dto.ConvertUserToResponse(user)

	res = dto.Response{
		Data:   &userRes,
		Mess:   "ok",
		Status: 200,
	}
	render.JSON(w, r, res)
}

func (c controller) FilterUser(w http.ResponseWriter, r *http.Request) {
	var page, pageSize int
	var err error
	var res dto.ListResponse

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

	searchFilter := dto.GetFilterSearch(search, status, sortBy, direction, page, pageSize)
	users, total, err := c.UserService.FilterUser(searchFilter)
	if err != nil {
		res = dto.ListResponse{
			Mess:     err.Error(),
			Status:   500,
			Data:     nil,
			Page:     page,
			PageSize: pageSize,
			Total:    0,
		}

		render.JSON(w, r, res)
		return
	}

	userRes := make([]*dto.ResponseUser, 0)
	for _, user := range users {
		ele := dto.ConvertUserToResponse(user)
		userRes = append(userRes, &ele)
	}

	res = dto.ListResponse{
		Mess:     "ok",
		Status:   200,
		Data:     userRes,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	render.JSON(w, r, res)
}

func (c controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var res dto.Response
	var user dto.UpdateUser

	id := chi.URLParam(r, "id")
	if id == "" {
		res = dto.Response{
			Data:   nil,
			Mess:   "id is required",
			Status: 0,
		}

		render.JSON(w, r, res)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	updateUser, err := c.UserService.UpdateUser(user, id)
	if err != nil {
		res = dto.Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)

		return
	}
	userRes := dto.ConvertUserToResponse(updateUser)
	res = dto.Response{
		Data:   &userRes,
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
	Register(w http.ResponseWriter, r *http.Request)
}

func NewController(client *mongo.Client) Controller {
	return controller{
		UserService: NewService(client, infrastructure.UserCollection, infrastructure.DatabaseName),
		RoleService: role.NewService(client, infrastructure.RoleCollection, infrastructure.DatabaseName),
	}
}
