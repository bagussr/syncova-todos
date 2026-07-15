package domain

type BasePaginationRequest struct {
	Page    int    `json:"page" form:"page" default:"1"`
	PerPage int    `json:"per_page" form:"per_page" default:"10"`
	SortBy  string `json:"sort_by" form:"sort_by" default:"created_at"`
	Sort    string `json:"sort" form:"sort" default:"desc"`
	Search  string `json:"search" form:"search"`
}
