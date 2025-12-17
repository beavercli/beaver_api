package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/beavercli/beaver_api/internal/integrations/github"
	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
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

	accessToken, err := s.generateAccessToken(userID, 24*7*time.Hour)
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
				RefreshToken: "", // TODO
			},
		}}, nil
}

type AccessTokenClaims struct {
	jwt.Claims
	UserID int64  `json:"uid"`
	Scope  string `json:"scope,omitempty"`
}

// TODO: move lthe JWT/JWE logic to a separe pkg
func (s *Service) generateAccessToken(userID int64, ttl time.Duration) (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: s.conf.Secret},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := AccessTokenClaims{
		Claims: jwt.Claims{
			Subject:   "user",
			Issuer:    "beaver_api",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Expiry:    jwt.NewNumericDate(now.Add(ttl)),
		},
		UserID: userID,
		Scope:  "api",
	}

	return jwt.Signed(signer).Claims(claims).Serialize()
}

func generateRandomPwd() (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func (s *Service) generateJWE(g github.GithubDevicePayload) (string, error) {
	r := jose.Recipient{
		Key:       s.conf.Secret,
		Algorithm: jose.DIRECT,
	}
	opts := jose.EncrypterOptions{}
	opts.WithContentType("JWT")

	enc, err := jose.NewEncrypter(jose.A256GCM, r, &opts)
	if err != nil {
		return "", err
	}

	token, err := jwt.Encrypted(enc).Claims(toOAuthGithubJWE(g)).Serialize()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) decryptJWE(token string) (OAuthGithubJWE, error) {
	enc, err := jwt.ParseEncrypted(token, []jose.KeyAlgorithm{jose.DIRECT}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		return OAuthGithubJWE{}, err
	}
	var claims OAuthGithubJWE
	if err := enc.Claims(s.conf.Secret, &claims); err != nil {
		return OAuthGithubJWE{}, err
	}
	return claims, nil
}

func toOAuthGithubJWE(g github.GithubDevicePayload) OAuthGithubJWE {
	return OAuthGithubJWE{
		DeviceCode: g.DeviceCode,
		ExpiresIn:  time.Now().Add(time.Duration(g.ExpiresIn) * time.Second).Unix(),
	}
}
