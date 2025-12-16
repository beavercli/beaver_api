package router

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

// Type aliases for Swagger documentation
type SnippetsPageResponse = PageResponse[SnippetSummary]
type TagsPageResponse = PageResponse[Tag]
type LanguagesPageResponse = PageResponse[Language]
type ContributorsPageResponse = PageResponse[Contributor]
