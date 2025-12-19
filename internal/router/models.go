package router

type RefreshToken struct {
	UserID       string `json:"user_id"` // TODO: REMOVE (only for test)
	RefreshToken string `json:"refresh_token"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Contributor struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Git struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Snippet struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Code         string        `json:"code"`
	Git          Git           `json:"git"`
	GitPath      string        `json:"git_path"`
	GitVersion   string        `json:"git_version"`
	ProjectURL   string        `json:"project_url,omitempty"`
	Language     Language      `json:"language"`
	Tags         []Tag         `json:"tags"`
	Contributors []Contributor `json:"contributors"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type SnippetSummary struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	ProjectURL string   `json:"project_url,omitempty"`
	Git        Git      `json:"git"`
	GitPath    string   `json:"git_path"`
	GitVersion string   `json:"git_version"`
	Language   Language `json:"language"`
	Tags       []Tag    `json:"tags"`
}

type SnippetListFilterArg struct {
	LanguageID *int64  // nil or >0
	TagIDs     []int64 // nil or all(>0)
}

type DeviceOAuth struct {
	UserCode  string `json:"user_code"`
	URL       string `json:"url"`
	Token     string `json:"token"`
	ExpiersIn int    `json:"expiers_in"`
	Interval  int    `json:"interval"`
}

type CreateContributorRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

type CreateLanguageRequest struct {
	Name string `json:"name"`
}

type CreateGit struct {
	URL string `json:"name"`
}

type IngestSnippetRequest struct {
	Title        string                     `json:"title"`
	Code         string                     `json:"code"`
	ProjectURL   string                     `json:"project_url,omitempty"`
	Git          CreateGit                  `json:"git_repo_url"`
	GitPath      string                     `json:"git_path"`
	GitVersion   string                     `json:"git_version"`
	Language     CreateLanguageRequest      `json:"language"`
	Tags         []CreateTagRequest         `json:"tags"`
	Contributors []CreateContributorRequest `json:"contributors"`
}

type GithubPullRequest struct {
	Token string `json:"token"`
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
	Status  DeviceAuthStatus `json:"status"`
	Session *Session
}

type CreateServiceAccessTokenRequest struct {
	// Human readable label to identify the token.
	Name string `json:"name"`
	// Optional ISO8601 expiry timestamp; omit for long-lived tokens.
	ExpiresAt string `json:"expires_at,omitempty"`
}

type ServiceAccessToken struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Token     string `json:"token,omitempty"` // Secret token value; returned only at creation time
	ExpiresAt string `json:"expires_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ServiceAccessTokenSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

// Type aliases for Swagger documentation
type SnippetsPageResponse = PageResponse[SnippetSummary]
type TagsPageResponse = PageResponse[Tag]
type LanguagesPageResponse = PageResponse[Language]
type ContributorsPageResponse = PageResponse[Contributor]
type ServiceAccessTokensPageResponse = PageResponse[ServiceAccessTokenSummary]
