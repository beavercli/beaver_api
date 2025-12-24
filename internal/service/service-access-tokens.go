package service

import (
	"context"
	"time"
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

	// TODO: save hash of token

	return ServiceAccessToken{
		Name:      args.Name,
		ExpiresAt: args.ExpiresAt,
		Token:     t,
	}, nil
}
