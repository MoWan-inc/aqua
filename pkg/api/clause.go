package api

type Pagination struct {
	// 当前页码, 从1开始
	Page int `form:"page" json:"page" validate:"gte=1"`
	// 每页数量
	PageSize int `form:"page_size" json:"page_size" validate:"gte=1,lte=100"`
}

func (p *Pagination) Empty() bool {
	return p.Page <= 0 || p.PageSize == 0
}

type Sorting struct {
	SortBy   string `form:"sort_by" json:"sort_by"`
	SortDesc bool   `form:"sort_desc" json:"sort_desc"`
}

type Filter struct {
	Filters string `form:"filters" json:"filters"`
	Fields  string `form:"fields" json:"fields"`
}

// PaginationSortingFilter 分页、排序、过滤, 用于查询请求的过滤条件
type PaginationSortingFilter struct {
	Pagination
	Sorting
	Filter
}
