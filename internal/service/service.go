package service

import "github.com/beavercli/beaver_api/internal/storage"

type Service struct {
	db *storage.Queries
}

func New(db *storage.Queries) *Service {
	return &Service{db: db}
}
