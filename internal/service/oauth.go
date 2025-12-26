package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/beavercli/beaver_api/internal/integrations/github"
	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetDeviceRequest(ctx context.Context) (OAuthRedirect, error) {
	g, err := s.github.GetDeviceCode(ctx)

	token, err := s.generateJWE(g)
	if err != nil {
		return OAuthRedirect{}, err
	}

	return OAuthRedirect{
		URL:       g.VerificationURL,
		UserCode:  g.UserCode,
		ExpiresIn: g.ExpiresIn,
		Interval:  g.Interval,
		Token:     token,
	}, nil
}

func (s *Service) GithubDevicePoll(ctx context.Context, jwe string) (DeviceAuthResult, error) {
	dc, err := s.decryptJWE(jwe)
	if err != nil {
		return DeviceAuthResult{}, err
	}

	if time.Now().Unix() > dc.ExpiresIn {
		return DeviceAuthResult{Status: DeviceAuthExpired, Session: nil}, nil
	}

	at, err := s.github.GetAccessToken(ctx, dc.DeviceCode)
	if err != nil {
		return DeviceAuthResult{}, err
	}
	if at.Error != "" {
		return DeviceAuthResult{Status: DeviceAuthPending, Session: nil}, nil
	}

	g, ctxG := errgroup.WithContext(ctx)

	var githubUser github.GithubUserPayload
	g.Go(func() error {
		var err error
		githubUser, err = s.github.GetUser(ctxG, at)
		return err
	})

	var githubUserEmail github.GithubUserEmailPayload
	g.Go(func() error {
		var err error
		githubUserEmail, err = s.github.GetUserEmail(ctxG, at)
		return err
	})

	if err := g.Wait(); err != nil {
		return DeviceAuthResult{}, err
	}

	pwd, err := generateRandomPwd()
	if err != nil {
		return DeviceAuthResult{}, err
	}

	userID, err := s.db.UpsertUser(ctx, storage.UpsertUserParams{
		Username:     githubUser.Login,
		Email:        githubUserEmail.Email,
		PasswordHash: pwd,
	})
	if err != nil {
		return DeviceAuthResult{}, err
	}

	accessToken, err := s.IssueJWT(AccessToken, userID, AccessTokenTTL)
	if err != nil {
		return DeviceAuthResult{}, err
	}

	refreshToken, err := s.IssueJWT(RefreshToken, userID, RefreshTokenTTL)

	_, err = s.db.CreateRefreshToken(ctx, storage.CreateRefreshTokenParams{
		UserID:    pgtype.Int8{Int64: userID, Valid: true},
		TokenHash: computeHash(refreshToken),
		IssuedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(RefreshTokenTTL), Valid: true},
	})
	if err != nil {
		return DeviceAuthResult{}, err
	}

	return DeviceAuthResult{
		Status: DeviceAuthDone,
		Session: &Session{
			User: User{
				ID:       userID,
				Email:    githubUserEmail.Email,
				Username: githubUser.Login,
			},
			TokenPair: TokenPair{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}}, nil
}

func (s *Service) RotateTokens(ctx context.Context, userID int64, refreshToken string) (TokenPair, error) {
	c, err := s.ParseJWT(refreshToken)
	if err != nil {
		return TokenPair{}, err
	}
	// check ExpiredAt in the JWT token
	tn := time.Now().Unix()
	if c.Expiry.Time().Unix() < tn {
		return TokenPair{}, fmt.Errorf("Refresh token expired")
	}

	// check the token exist in the issued refresh tokens
	t, err := s.db.GetRefreshTokenByHash(ctx, computeHash(refreshToken))
	if err != nil {
		return TokenPair{}, err
	}
	// there is an edge case when we need to update the TTL for the
	// refresh tokens. For that case we also need to check ExpiredAt
	// with the value stored in the DB
	if t.ExpiresAt.Time.Unix() < tn {
		return TokenPair{}, fmt.Errorf("Refresh token expired")
	}

	// all is good we can issue a new token pair

	at, err := s.IssueJWT(AccessToken, userID, AccessTokenTTL)
	if err != nil {
		return TokenPair{}, err
	}

	rt, err := s.IssueJWT(RefreshToken, userID, RefreshTokenTTL)
	if err != nil {
		return TokenPair{}, err
	}

	txOpts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	err = s.inTx(ctx, txOpts, func(db *storage.Queries) error {
		if err := s.db.DeleteRefreshTokenByID(ctx, t.ID); err != nil {
			return err
		}

		if _, err := s.db.CreateRefreshToken(ctx, storage.CreateRefreshTokenParams{
			UserID:    pgtype.Int8{Int64: userID, Valid: true},
			TokenHash: computeHash(rt),
			IssuedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
			ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(RefreshTokenTTL), Valid: true},
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) AuthUser(ctx context.Context, tokenType TokenType, token string) (int64, error) {
	switch tokenType {
	case AccessToken:
		return s.handleAccessToken(ctx, token)
	case SessionToken:
		return s.handleSessionToken(ctx, token)
	default:
		return 0, fmt.Errorf("Provided token is not supported yet: %s", tokenType)
	}
}

func (s *Service) handleAccessToken(_ context.Context, token string) (int64, error) {
	c, err := s.ParseJWT(token)
	if err != nil {
		return 0, err
	}

	if time.Now().After(c.Expiry.Time()) {
		return 0, fmt.Errorf("Token is expired")
	}

	userID, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
func (s *Service) handleSessionToken(ctx context.Context, token string) (int64, error) {
	c, err := s.ParseJWT(token)
	if err != nil {
		return 0, err
	}

	if time.Now().After(c.Expiry.Time()) {
		return 0, fmt.Errorf("Session token is expired")
	}

	userID, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	t, err := s.db.GetRefreshTokenByHash(ctx, computeHash(token))
	if err != nil {
		return 0, err
	}

	if time.Now().After(t.ExpiresAt.Time) {
		return 0, fmt.Errorf("Session token is expired based on the udpated expiry")
	}

	return userID, nil

}

func computeHash(t string) string {
	sum := sha256.Sum256([]byte(t))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func generateRandomPwd() (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}
