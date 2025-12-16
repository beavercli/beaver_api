package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	githubDeviceUrl = "https://github.com/login/device/code"
	scope           = "read:user user:email"
)

type GithubDeviceRequest struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURL string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func (s *Service) GetDeviceRequest(ctx context.Context) (OAuthRedirect, error) {
	val := url.Values{
		"client_id": {s.oauthConf.ClinetID},
		"scope":     {scope},
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rq, err := http.NewRequestWithContext(ctx, "POST", githubDeviceUrl, strings.NewReader(val.Encode()))
	if err != nil {
		return OAuthRedirect{}, err
	}
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return OAuthRedirect{}, err
	}
	defer rp.Body.Close()

	if rp.StatusCode != http.StatusOK {
		fmt.Println(rp.Body) // todo: remove
		return OAuthRedirect{}, fmt.Errorf("Request device code failed: %s", rp.Status)
	}

	var rb GithubDeviceRequest
	if err := json.NewDecoder(rp.Body).Decode(&rb); err != nil {
		return OAuthRedirect{}, err
	}

	token, err := s.generateJWE(rb)
	if err != nil {
		return OAuthRedirect{}, err
	}

	return OAuthRedirect{
		URL:       rb.VerificationURL,
		UserCode:  rb.UserCode,
		ExpiresIn: rb.ExpiresIn,
		Interval:  rb.Interval,
		Token:     token,
	}, nil
}

func (s *Service) generateJWE(g GithubDeviceRequest) (string, error) {
	r := jose.Recipient{
		Key:       s.oauthConf.Secret,
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
func toOAuthGithubJWE(g GithubDeviceRequest) OAuthGithubJWE {
	return OAuthGithubJWE{
		DeviceCode: g.DeviceCode,
		ExpiresIn:  g.ExpiresIn,
	}
}
