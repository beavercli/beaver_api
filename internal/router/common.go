package router

import (
	"fmt"
	"net/url"
	"strconv"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type PageQueryArg struct {
	Page     int // default 1
	PageSize int // default 20, max 100
}

func toPageQuery(v url.Values) (PageQueryArg, error) {
	page := defaultPage
	pageSize := defaultPageSize

	if raw := v.Get("page"); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil || val <= 0 {
			return PageQueryArg{}, fmt.Errorf("page must be a positive integer")
		}
		page = val
	}
	if raw := v.Get("page_size"); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil || val < 1 {
			return PageQueryArg{}, fmt.Errorf("page_size must be a positive integer")
		}
		if val > maxPageSize {
			return PageQueryArg{}, fmt.Errorf("page_size mst be <=%d", maxPageSize)
		}
		pageSize = val
	}

	return PageQueryArg{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

type PageResponse[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

func toPage[T any](items []T, total int, page, pageSize int) PageResponse[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return PageResponse[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
