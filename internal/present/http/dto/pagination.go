package dto

import "strings"

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type PaginationRequest struct {
	Page     int    `json:"page" form:"page" query:"page"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" form:"sort_by" query:"sort_by"`
	Order    string `json:"order" form:"order" query:"order"` // "asc" or "desc"
}

func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	p.Order = strings.ToLower(strings.TrimSpace(p.Order))
	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "asc"
	}
}

func (p PaginationRequest) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p PaginationRequest) Limit() int {
	if p.PageSize <= 0 {
		return DefaultPageSize
	}
	return p.PageSize
}

type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationResponse[T any](items []T, total int64, req PaginationRequest) PaginationResponse[T] {
	req.Normalize()
	var totalPages int
	if total > 0 {
		totalPages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}
	return PaginationResponse[T]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}
}
