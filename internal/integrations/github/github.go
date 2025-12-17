package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GithubDeviceUrl    = "https://github.com/login/device/code"
	GithubAcessToken   = "https://github.com/login/oauth/access_token"
	GetGithubUser      = "https://api.github.com/user"
	GetGithubUserEmail = "https://api.github.com/user/emails"
	Scope              = "read:user user:email"
	GrantType          = "urn:ietf:params:oauth:grant-type:device_code"
)

type Client struct {
	Timeout  time.Duration
	ClientID string
}

func New(clientID string, timeout time.Duration) *Client {
	return &Client{
		Timeout:  timeout,
		ClientID: clientID,
	}
}

type GithubDevicePayload struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURL string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func (c *Client) GetDeviceCode(ctx context.Context) (GithubDevicePayload, error) {
	val := url.Values{
		"client_id": {c.ClientID},
		"scope":     {Scope},
	}
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	rq, err := http.NewRequestWithContext(ctx, "POST", GithubDeviceUrl, strings.NewReader(val.Encode()))
	if err != nil {
		return GithubDevicePayload{}, err
	}
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return GithubDevicePayload{}, err
	}
	defer rp.Body.Close()

	if rp.StatusCode != http.StatusOK {
		return GithubDevicePayload{}, fmt.Errorf("Request device code failed: %s", rp.Status)
	}

	var rb GithubDevicePayload
	if err := json.NewDecoder(rp.Body).Decode(&rb); err != nil {
		return GithubDevicePayload{}, err
	}
	return rb, nil
}

type GithubAccesTokenPayload struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`

	// error payload
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

func (c *Client) GetAccessToken(ctx context.Context, dc string) (GithubAccesTokenPayload, error) {
	val := url.Values{
		"client_id":   {c.ClientID},
		"scope":       {Scope},
		"grant_type":  {GrantType},
		"device_code": {dc},
	}
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	rq, err := http.NewRequestWithContext(ctx, "POST", GithubAcessToken, strings.NewReader(val.Encode()))
	if err != nil {
		return GithubAccesTokenPayload{}, err
	}
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return GithubAccesTokenPayload{}, err
	}
	defer rp.Body.Close()

	if rp.StatusCode != http.StatusOK {
		return GithubAccesTokenPayload{}, fmt.Errorf("Request access token failed: %s", rp.Status)
	}

	var rb GithubAccesTokenPayload
	if err := json.NewDecoder(rp.Body).Decode(&rb); err != nil {
		return GithubAccesTokenPayload{}, err
	}
	return rb, nil
}

type GithubUserPayload struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

func (c *Client) GetUser(ctx context.Context, t GithubAccesTokenPayload) (GithubUserPayload, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	rq, err := http.NewRequestWithContext(ctx, "GET", GetGithubUser, nil)
	if err != nil {
		return GithubUserPayload{}, err
	}
	rq.Header.Add("Accept", "application/vnd.github+json")
	rq.Header.Add("Authorization", "Bearer "+t.AccessToken)

	rp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return GithubUserPayload{}, err
	}
	defer rp.Body.Close()

	if rp.StatusCode != http.StatusOK {
		return GithubUserPayload{}, fmt.Errorf("Request github user failed: %s", rp.Status)
	}

	var rb GithubUserPayload
	if err := json.NewDecoder(rp.Body).Decode(&rb); err != nil {
		return GithubUserPayload{}, err
	}
	return rb, nil
}

type GithubUserEmailPayload struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func (c *Client) GetUserEmail(ctx context.Context, t GithubAccesTokenPayload) (GithubUserEmailPayload, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	rq, err := http.NewRequestWithContext(ctx, "GET", GetGithubUserEmail, nil)
	if err != nil {
		return GithubUserEmailPayload{}, err
	}
	rq.Header.Add("Accept", "application/vnd.github+json")
	rq.Header.Add("Authorization", "Bearer "+t.AccessToken)

	rp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return GithubUserEmailPayload{}, err
	}
	defer rp.Body.Close()

	if rp.StatusCode != http.StatusOK {
		return GithubUserEmailPayload{}, fmt.Errorf("Request github user failed: %s", rp.Status)
	}

	var rb []GithubUserEmailPayload
	if err := json.NewDecoder(rp.Body).Decode(&rb); err != nil {
		return GithubUserEmailPayload{}, err
	}
	var email GithubUserEmailPayload
	for _, e := range rb {
		if e.Primary {
			email = e
			break
		}
	}

	// fallback if for some reason there is not a primary email
	if email.Email == "" {
		email = rb[0]
	}
	return email, nil
}
