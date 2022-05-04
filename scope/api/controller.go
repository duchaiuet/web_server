package api

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

func (c controller) CreateApi(w http.ResponseWriter, r *http.Request) {
	var res Response
	var api CUAPI

	err := json.NewDecoder(r.Body).Decode(&api)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	createApi, err := c.Service.CreateAPI(api)
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
		Data:   createApi,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) GetApiById(w http.ResponseWriter, r *http.Request) {
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

	api, err := c.Service.GetAPIById(id)
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
		Data:   api,
		Mess:   "ok",
		Status: 200,
	}
	render.JSON(w, r, res)
}

func (c controller) FilterApi(w http.ResponseWriter, r *http.Request) {
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
	apis, err := c.Service.FilterAPI(searchFilter)
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
		Data:     apis,
		Page:     page,
		PageSize: pageSize,
	}

	render.JSON(w, r, res)
}

func (c controller) UpdateApi(w http.ResponseWriter, r *http.Request) {
	var res Response
	var api CUAPI

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

	err := json.NewDecoder(r.Body).Decode(&api)
	if err != nil {
		res = Response{
			Data:   nil,
			Mess:   err.Error(),
			Status: 400,
		}
		render.JSON(w, r, res)
		return
	}

	updateApi, err := c.Service.UpdateAPI(api, id)
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
		Data:   updateApi,
		Mess:   "ok",
		Status: 200,
	}

	render.JSON(w, r, res)
}

func (c controller) DeleteApi(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

type Controller interface {
	CreateApi(w http.ResponseWriter, r *http.Request)
	GetApiById(w http.ResponseWriter, r *http.Request)
	FilterApi(w http.ResponseWriter, r *http.Request)
	UpdateApi(w http.ResponseWriter, r *http.Request)
	DeleteApi(w http.ResponseWriter, r *http.Request)
}

func NewController(client *mongo.Client) Controller {
	return controller{
		Service: NewService(client, infrastructure.APICollection, infrastructure.DatabaseName),
	}
}
