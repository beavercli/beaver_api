package service

import "github.com/beavercli/beaver_api/internal/storage"

type Service struct {
	db *storage.Queries
}

func New(db *storage.Queries) *Service {
	return &Service{db: db}
}

type PageParam struct {
	Page     int
	PageSize int
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PageParam) Limit() int {
	return p.PageSize
}
