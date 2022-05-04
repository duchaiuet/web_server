package api_permission

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"web_server/infrastructure"
)

type controller struct {
	Service Service
}

func (c controller) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var res Response
	var rule CUCasbinRule

	err := json.NewDecoder(r.Body).Decode(&rule)

	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	createPermission, err := c.Service.CreateRule(rule)
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
		Data:   createPermission,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) GetPermissionById(w http.ResponseWriter, r *http.Request) {
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

	permission, err := c.Service.GetRuleById(id)
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
		Data:   permission,
		Mess:   "ok",
		Status: 200,
	}
	render.JSON(w, r, res)
}

func (c controller) FilterPermission(w http.ResponseWriter, r *http.Request) {
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
	permissions, err := c.Service.FilterRule(searchFilter)
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
		Data:     permissions,
		Page:     page,
		PageSize: pageSize,
	}

	render.JSON(w, r, res)
}

func (c controller) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	var res Response
	var rule CUCasbinRule

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

	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	updatePermission, err := c.Service.UpdateRule(rule, id)
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
		Data:   updatePermission,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) DeletePermission(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

type Controller interface {
	CreatePermission(w http.ResponseWriter, r *http.Request)
	GetPermissionById(w http.ResponseWriter, r *http.Request)
	FilterPermission(w http.ResponseWriter, r *http.Request)
	UpdatePermission(w http.ResponseWriter, r *http.Request)
	DeletePermission(w http.ResponseWriter, r *http.Request)
}

func NewController(client *mongo.Client) Controller {
	return controller{
		Service: NewService(client, infrastructure.CasbinRuleCollection, infrastructure.DatabaseName),
	}
}
