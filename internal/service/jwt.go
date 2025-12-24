package service

import (
	"strconv"
	"time"

	"github.com/beavercli/beaver_api/internal/integrations/github"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessTokenTTL  time.Duration = 24 * 7 * time.Hour
	RefreshTokenTTL time.Duration = AccessTokenTTL * 4
	AccessToken     TokenType     = "access"
	RefreshToken    TokenType     = "refresh"
	SessionToken    TokenType     = "session"
)

type JWTClaims struct {
	jwt.Claims
	Type TokenType `json:"token_type"`
}

func (s *Service) IssueJWT(tt TokenType, userID int64, ttl time.Duration) (string, error) {
	signerOpts := jose.SignerOptions{}
	signerOpts.WithType("JWT")

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: s.conf.Secret},
		&signerOpts,
	)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := JWTClaims{
		Claims: jwt.Claims{
			ID:        uuid.New().String(),
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    "beaver_api",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Expiry:    jwt.NewNumericDate(now.Add(ttl)),
		},
		Type: tt,
	}

	return jwt.Signed(signer).Claims(claims).Serialize()
}

func (s *Service) ParseJWT(t string) (JWTClaims, error) {
	token, err := jwt.ParseSigned(t, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return JWTClaims{}, err
	}
	claims := JWTClaims{}
	if err := token.Claims(s.conf.Secret, &claims); err != nil {
		return JWTClaims{}, err
	}
	return claims, nil
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
