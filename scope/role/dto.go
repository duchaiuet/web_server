package role

import "web_server/model"

type SearchFilter struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type CURole struct {
	Code   string `json:"code" bson:"code"`
	Name   string `json:"name" bson:"name"`
	Status bool   `json:"status" bson:"status"`
}

type Response struct {
	Data   *model.Role `json:"data"`
	Mess   string      `json:"mess"`
	Status int         `json:"status"`
}

type ListResponse struct {
	Mess     string        `json:"mess"`
	Status   int           `json:"status"`
	Data     []*model.Role `json:"data"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

func ConvertCURoleToModel(cuROle CURole) model.Role {
	return model.Role{
		Code:   cuROle.Code,
		Name:   cuROle.Name,
		Status: cuROle.Status,
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
