package user

import (
	"web_server/model"
)

type SearchFilter struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	SortBy   SortBy `json:"sort_by"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type SortBy struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}

type CreateUser struct {
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Status       string `json:"status"`
	Role         int    `json:"role"`
}

type UpdateUser struct {
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Status       string `json:"status"`
	Role         int    `json:"role"`
}

type Login struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Response struct {
	User   *model.User `json:"user"`
	Mess   string      `json:"mess"`
	Status int         `json:"status"`
}

type ListResponse struct {
	Mess     string        `json:"mess"`
	Status   int           `json:"status"`
	Users    []*model.User `json:"users"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type LoginResponse struct {
	User   *model.User `json:"user"`
	Token  string      `json:"token"`
	Status int         `json:"status"`
	Mess   string      `json:"mess"`
}

func ConvertCreteUserToModel(create CreateUser) model.User {
	return model.User{
		UserName:     create.UserName,
		Password:     create.Password,
		EmailAddress: create.EmailAddress,
		FirstName:    create.FirstName,
		LastName:     create.LastName,
		Status:       create.Status,
	}
}

func ConvertUpdateUserToModel(update UpdateUser) model.User {
	return model.User{
		UserName:     update.UserName,
		Password:     update.Password,
		EmailAddress: update.EmailAddress,
		FirstName:    update.FirstName,
		LastName:     update.LastName,
		Status:       update.Status,
	}
}

func GetFilterSearch(search string, status string, sortBy string, direction string, page int, pageSize int) SearchFilter {
	filter := SearchFilter{
		Search: search,
		Status: status,
		SortBy: SortBy{
			Field: sortBy,
			Sort:  direction,
		},
		Page:     page,
		PageSize: pageSize,
	}

	return filter
}
