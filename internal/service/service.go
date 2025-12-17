package service

import (
	"context"

	"github.com/beavercli/beaver_api/internal/integrations/github"
	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Secret []byte
}

type GithubOAuthClient interface {
	GetDeviceCode(ctx context.Context) (github.GithubDevicePayload, error)
	GetAccessToken(ctx context.Context, dc string) (github.GithubAccesTokenPayload, error)
	GetUser(ctx context.Context, t github.GithubAccesTokenPayload) (github.GithubUserPayload, error)
	GetUserEmail(ctx context.Context, t github.GithubAccesTokenPayload) (github.GithubUserEmailPayload, error)
}

type Service struct {
	conf   Config
	github GithubOAuthClient
	pool   *pgxpool.Pool
	db     *storage.Queries
}

func New(pool *pgxpool.Pool, c Config, github GithubOAuthClient) *Service {
	return &Service{
		conf:   c,
		github: github,
		pool:   pool,
		db:     storage.New(pool),
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
