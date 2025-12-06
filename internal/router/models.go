package router

import (
	"net/url"
	"strconv"

	"github.com/beavercli/beaver_api/internal/service"
)

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toTag(t service.Tag) Tag {
	return Tag{
		ID:   strconv.FormatInt(t.ID, 10),
		Name: t.Name,
	}
}

func toTags(ts []service.Tag) []Tag {
	tags := make([]Tag, len(ts))
	for i, t := range ts {
		tags[i] = toTag(t)
	}
	return tags
}

type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toLanguage(l *service.Language) *Language {
	if l == nil {
		return nil
	}

	return &Language{
		ID:   strconv.FormatInt(l.ID, 10),
		Name: l.Name,
	}
}

type Contributor struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func toContributor(c service.Contributor) Contributor {
	return Contributor{
		ID:        strconv.FormatInt(c.ID, 10),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
	}
}

func toContributors(cs []service.Contributor) []Contributor {
	contributors := make([]Contributor, len(cs))
	for i, c := range cs {
		contributors[i] = toContributor(c)
	}
	return contributors
}

type Snippet struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Code         string        `json:"code"`
	ProjectURL   string        `json:"project_url,omitempty"`
	Language     *Language     `json:"language,omitempty"`
	Tags         []Tag         `json:"tags"`
	Contributors []Contributor `json:"contributors"`
}

func toSnippet(s service.Snippet) Snippet {
	return Snippet{
		ID:           strconv.FormatInt(s.ID, 10),
		Title:        s.Title,
		Code:         s.Code,
		ProjectURL:   s.ProjectURL,
		Language:     toLanguage(s.Language),
		Tags:         toTags(s.Tags),
		Contributors: toContributors(s.Contributors),
	}
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type SnippetSummary struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	ProjectURL string    `json:"project_url,omitempty"`
	Language   *Language `json:"language,omitempty"`
	Tags       []Tag     `json:"tags"`
}

func toSnippetSummary(s service.SnippetSummary) SnippetSummary {
	return SnippetSummary{
		ID:         strconv.FormatInt(s.ID, 10),
		Title:      s.Title,
		ProjectURL: s.ProjectURL,
		Language:   toLanguage(s.Language),
		Tags:       toTags(s.Tags),
	}
}

func toSnippetSummaries(ss []service.SnippetSummary) []SnippetSummary {
	snippetSummers := make([]SnippetSummary, len(ss))
	for i, s := range ss {
		snippetSummers[i] = toSnippetSummary(s)
	}
	return snippetSummers
}

type SnippetListFilterArg struct {
	LanguageID *int64
	TagIDs     []int64
}

func toSnippetListFilterArg(v url.Values) SnippetListFilterArg {
	return SnippetListFilterArg{
		LanguageID: strToInt(v.Get("language_id"), nil),
		TagIDs:     strToInts(v["tag_id"]),
	}
}

type PageQueryArg struct {
	Page     int
	PageSize int
}

func toPageQuery(v url.Values) PageQueryArg {
	return PageQueryArg{
		Page:     int(*strToInt(v.Get("page"), intPtr(1))),
		PageSize: int(*strToInt(v.Get("page_size"), intPtr(20))),
	}
}

type PageResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func toPage[T any](items []T, total int64, page, pageSize int) PageResponse[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return PageResponse[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Type aliases for Swagger documentation
type SnippetsPageResponse = PageResponse[SnippetSummary]
type TagsPageResponse = PageResponse[Tag]
type LanguagesPageResponse = PageResponse[Language]
type ContributorsPageResponse = PageResponse[Contributor]

type CreateSnippetRequest struct {
	Title        string   `json:"title"`
	Code         string   `json:"code"`
	ProjectURL   string   `json:"project_url,omitempty"`
	Language     string   `json:"language"`
	Tags         []string `json:"tags"`
	Contributors []string `json:"contributors"`
}
