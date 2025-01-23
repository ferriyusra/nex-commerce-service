package response

type ErrorResponseDefault struct {
	Meta
}

type Meta struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type DefaultSucceesResponse struct {
	Meta       Meta                `json:"meta"`
	Data       interface{}         `json:"data,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	TotalRecords int `json:"totalRecrods"`
	Page         int `json:"page"`
	PerPage      int `json:"perPage"`
	TotalPages   int `json:"totalPages"`
}
