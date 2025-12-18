package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/beavercli/beaver_api/internal/integrations/github"
	"github.com/beavercli/beaver_api/internal/storage"
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

func (s *Service) UpsertUser(ctx context.Context, jwe string) (DeviceAuthResult, error) {
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
		TokenHash: hashRefreshToken(refreshToken),
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

func hashRefreshToken(t string) string {
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
