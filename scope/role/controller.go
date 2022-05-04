package role

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

type controller struct {
	Service Service
}

func (c controller) CreateRole(w http.ResponseWriter, r *http.Request) {
	var res Response
	var role CURole

	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	createRole, err := c.Service.CreateRole(role)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)
		return
	}
	res = Response{
		Data:   createRole,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) GetRole(w http.ResponseWriter, r *http.Request) {
	var res Response
	id := chi.URLParam(r, "id")
	if id == "" {
		res = Response{
			Data:   nil,
			Mess:   "id is required",
			Status: 0,
		}

		render.JSON(w, r, res)
		return
	}

	role, err := c.Service.GetRoleById(id)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}

		render.JSON(w, r, res)
		return
	}
	res = Response{
		Data:   role,
		Mess:   "ok",
		Status: 200,
	}
	render.JSON(w, r, res)
}

func (c controller) FilterRole(w http.ResponseWriter, r *http.Request) {
	var page, pageSize int
	var err error
	var res ListResponse

	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	status := r.URL.Query().Get("status")

	if page, err = strconv.Atoi(pageStr); err == nil {
		page = 1
	}

	if pageSize, err = strconv.Atoi(pageSizeStr); err == nil {
		pageSize = 10
	}

	searchFilter := GetFilterSearch(search, status, page, pageSize)
	users, err := c.Service.FilterRole(searchFilter)
	if err != nil {
		res = ListResponse{
			Mess:     err.Error(),
			Status:   500,
			Data:     nil,
			Page:     page,
			PageSize: pageSize,
		}

		render.JSON(w, r, res)
		return
	}

	res = ListResponse{
		Mess:     "ok",
		Status:   200,
		Data:     users,
		Page:     page,
		PageSize: pageSize,
	}

	render.JSON(w, r, res)
}

func (c controller) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var res Response
	var user CURole

	id := chi.URLParam(r, "id")
	if id == "" {
		res = Response{
			Data:   nil,
			Mess:   "id is required",
			Status: 0,
		}

		render.JSON(w, r, res)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	updateRole, err := c.Service.UpdateRole(user, id)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 500,
		}
		render.JSON(w, r, res)

		return
	}

	res = Response{
		Data:   updateRole,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) DeleteRole(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

type Controller interface {
	CreateRole(w http.ResponseWriter, r *http.Request)
	GetRole(w http.ResponseWriter, r *http.Request)
	FilterRole(w http.ResponseWriter, r *http.Request)
	UpdateRole(w http.ResponseWriter, r *http.Request)
	DeleteRole(w http.ResponseWriter, r *http.Request)
}

func NewController(client *mongo.Client, collection string, database string) Controller {
	return controller{
		Service: NewService(client, collection, database),
	}
}
