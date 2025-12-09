package service

import (
	"context"

	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

type Service struct {
	pool *pgxpool.Pool
	db   *storage.Queries
}

func New(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
		db:   storage.New(pool),
	}
}

func (s *Service) inTx(ctx context.Context, txOpts pgx.TxOptions, fn func(q *storage.Queries) error) error {
	tx, err := s.pool.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}
	qtx := s.db.WithTx(tx)
	if err := fn(qtx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
