package api

import "web_server/model"

type SearchFilter struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type CUAPI struct {
	Path   string `json:"path"`
	Status bool   `json:"status"`
}

type Response struct {
	Data   *model.Api `json:"data"`
	Mess   string     `json:"mess"`
	Status int        `json:"status"`
}

type ListResponse struct {
	Mess     string       `json:"mess"`
	Status   int          `json:"status"`
	Data     []*model.Api `json:"data"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

func ConvertCUAPIToModel(cuAPI CUAPI) model.Api {
	return model.Api{
		Path:   cuAPI.Path,
		Status: cuAPI.Status,
	}
}

func GetFilterSearch(search string, status string, page int, pageSize int) SearchFilter {
	filter := SearchFilter{
		Search:   search,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}

	return filter
}
