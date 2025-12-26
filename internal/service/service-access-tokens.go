package service

import (
	"context"
	"fmt"
	"time"

	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
)

type CreateServiceAccessTokenArgs struct {
	UserID    int64
	Name      string
	ExpiresAt time.Duration
}

func (s *Service) CreateServceAccessToken(ctx context.Context, args CreateServiceAccessTokenArgs) (ServiceAccessToken, error) {
	t, err := s.IssueJWT(SessionToken, args.UserID, args.ExpiresAt)
	if err != nil {
		return ServiceAccessToken{}, err
	}
	arg := storage.CreateServiceAccessTokenParams{
		Name:      args.Name,
		TokenHash: computeHash(t),
		UserID:    pgtype.Int8{Int64: args.UserID, Valid: true},
		IssuedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(args.ExpiresAt), Valid: true},
	}
	fmt.Println(arg)
	st, err := s.db.CreateServiceAccessToken(ctx, arg)
	if err != nil {
		return ServiceAccessToken{}, err
	}

	return ServiceAccessToken{
		ID:        st.ID,
		Name:      st.Name,
		ExpiresAt: st.ExpiresAt.Time,
		IssuedAT:  st.IssuedAt.Time,
		Token:     t,
	}, nil
}

func (s *Service) ListServiceAccessTokens(ctx context.Context, userID int64, page PageParam) (ServiceAccessTokenList, error) {
	var at []storage.ServiceAccessToken
	var cnt int64

	g, dbCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		at, err = s.db.ListServiceAccessTokensByUserID(dbCtx, pgtype.Int8{Int64: userID, Valid: true})
		return err
	})
	g.Go(func() error {
		var err error
		cnt, err = s.db.CountServiceAccessTokensByUserID(dbCtx, pgtype.Int8{Int64: userID, Valid: true})
		return err
	})

	if err := g.Wait(); err != nil {
		return ServiceAccessTokenList{}, err
	}

	return ServiceAccessTokenList{
		Total: int(cnt),
		Items: toServiceAccessTokenSum(at),
	}, nil
}

func toServiceAccessTokenSum(sts []storage.ServiceAccessToken) []ServiceAccessTokenSum {
	s := make([]ServiceAccessTokenSum, len(sts))
	for i, st := range sts {
		s[i] = ServiceAccessTokenSum{
			ID:        st.ID,
			Name:      st.Name,
			ExpiresAt: st.ExpiresAt.Time,
			IssuedAT:  st.IssuedAt.Time,
		}
	}
	return s
}
