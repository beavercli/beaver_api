package service

import "time"

type Tag struct {
	ID   int64
	Name string
}

type TagList struct {
	Items []Tag
	Total int
}

type Language struct {
	ID   int64
	Name string
}

type LanguageList struct {
	Items []Language
	Total int
}

type Git struct {
	ID  int64
	URL string
}

type Contributor struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
}
type ContributorList struct {
	Items []Contributor
	Total int
}

type Snippet struct {
	ID           int64
	Title        string
	Code         string
	ProjectURL   string
	GitPath      string
	GitVersion   string
	Git          Git
	Language     Language
	Tags         []Tag
	Contributors []Contributor
}

type SnippetSummary struct {
	ID         int64
	Title      string
	ProjectURL string
	GitPath    string
	GitVersion string
	Git        Git
	Language   Language
	Tags       []Tag
}

type SnippetsList struct {
	Items []SnippetSummary
	Total int
}

type User struct {
	ID       int64
	Username string
	Email    string
}

type OAuthRedirect struct {
	UserCode  string
	URL       string
	Token     string
	ExpiresIn int
	Interval  int
}
type OAuthGithubJWE struct {
	DeviceCode string
	ExpiresIn  int64
}

type TokenPair struct {
	AccessToken  string `json:"acess_token"`
	RefreshToken string `json:"refresh_token"`
}

type Session struct {
	User      User
	TokenPair TokenPair
}

type DeviceAuthStatus string

const (
	DeviceAuthPending DeviceAuthStatus = "pending"
	DeviceAuthDone    DeviceAuthStatus = "done"
	DeviceAuthExpired DeviceAuthStatus = "expired"
)

type DeviceAuthResult struct {
	Status  DeviceAuthStatus
	Session *Session
}

type ServiceAccessToken struct {
	Name      string
	Token     string
	ExpiresAt time.Duration
}
