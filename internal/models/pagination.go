package models

type Pagination struct {
	PageSize     int    `json:"page_size"`
	PageNumber   int    `json:"page_number"`
	TotalCount   int    `json:"total_count"`
	TotalPages   int    `json:"total_pages"`
	NextPage     int    `json:"next_page,omitempty"`
	PreviousPage int    `json:"previous_page,omitempty"`
	SortBy       string `json:"sort_by,omitempty"`
	SortOrder    string `json:"sort_order,omitempty"`
	Data         any    `json:"data"`
}
