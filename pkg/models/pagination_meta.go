package models

type PaginationMeta struct {
	Page          int  `json:"page"`
	PageSize      int  `json:"page_size"`
	TotalElements int  `json:"total_elements"`
	TotalPages    int  `json:"total_pages"`
	IsFirst       bool `json:"is_first"`
	IsLast        bool `json:"is_last"`
}
