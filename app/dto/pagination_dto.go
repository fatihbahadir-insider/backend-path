package dto

type PaginationRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *PaginationRequest) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func NewPaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func NewPaginatedResponse[T any](data []T, page, limit int, total int64) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Data:       data,
		Pagination: NewPaginationMeta(page, limit, total),
	}
}