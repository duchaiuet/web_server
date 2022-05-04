package APIPermission

type SearchFilter struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
