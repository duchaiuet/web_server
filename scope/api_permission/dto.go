package api_permission

import (
	"web_server/model"
)

type SearchFilter struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type CUCasbinRule struct {
	ApiPath string `json:"api_path"`
	Role    string `json:"role"`
	Rule    string `json:"rule"`
}

type Response struct {
	Data   *model.CasbinRule `json:"data"`
	Mess   string            `json:"mess"`
	Status int               `json:"status"`
}

type ListResponse struct {
	Mess     string              `json:"mess"`
	Status   int                 `json:"status"`
	Data     []*model.CasbinRule `json:"data"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

func ConvertCUCasbinRuleToModel(cuAPI CUCasbinRule) model.CasbinRule {
	return model.CasbinRule{
		Path: cuAPI.ApiPath,
		Rule: cuAPI.Rule,
		Role: cuAPI.Role,
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
