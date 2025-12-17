package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/beavercli/beaver_api/internal/service"
)

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func toSnippet(s service.Snippet) Snippet {
	return Snippet{
		ID:           strconv.FormatInt(s.ID, 10),
		Title:        s.Title,
		Code:         s.Code,
		Git:          toGit(s.Git),
		GitPath:      s.GitPath,
		GitVersion:   s.GitVersion,
		ProjectURL:   s.ProjectURL,
		Language:     toLanguage(s.Language),
		Tags:         toTags(s.Tags),
		Contributors: toContributors(s.Contributors),
	}
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

func toLanguage(l service.Language) Language {
	return Language{
		ID:   strconv.FormatInt(l.ID, 10),
		Name: l.Name,
	}
}

func toGit(l service.Git) Git {
	return Git{
		ID:  strconv.FormatInt(l.ID, 10),
		URL: l.URL,
	}
}

func toLanguages(ls []service.Language) []Language {
	langs := make([]Language, len(ls))
	for i, l := range ls {
		langs[i] = toLanguage(l)
	}
	return langs
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

func toSnippetSummary(s service.SnippetSummary) SnippetSummary {
	return SnippetSummary{
		ID:         strconv.FormatInt(s.ID, 10),
		Title:      s.Title,
		ProjectURL: s.ProjectURL,
		GitPath:    s.GitPath,
		GitVersion: s.GitVersion,
		Git:        toGit(s.Git),
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

func toSnippetListFilterArg(v url.Values) (SnippetListFilterArg, error) {
	var langID *int64
	if raw := v.Get("language_id"); raw != "" {
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return SnippetListFilterArg{}, fmt.Errorf("language_id: %w", err)
		}
		if val <= 0 {
			return SnippetListFilterArg{}, fmt.Errorf("language_id must be positive")
		}
		langID = &val
	}

	tagParams := v["tag_id"]
	tags := make([]int64, 0, len(tagParams))
	for _, raw := range tagParams {
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return SnippetListFilterArg{}, fmt.Errorf("tag_id %q: %w", raw, err)
		}
		if val <= 0 {
			return SnippetListFilterArg{}, fmt.Errorf("tag_id must be positive")
		}
		tags = append(tags, val)
	}

	return SnippetListFilterArg{
		LanguageID: langID,
		TagIDs:     tags,
	}, nil
}

func toCreateSnippetRequestBody(r *http.Request) (IngestSnippetRequest, error) {
	defer r.Body.Close()

	var p IngestSnippetRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&p); err != nil {
		return IngestSnippetRequest{}, err
	}
	return p, nil
}

func toCreateSnippetParams(sr IngestSnippetRequest) service.CreateSnippetParam {
	ts := make([]service.CreateTagParam, len(sr.Tags))
	for i, t := range sr.Tags {
		ts[i] = service.CreateTagParam{
			Name: t.Name,
		}
	}

	cs := make([]service.CreateContributorParam, len(sr.Contributors))
	for i, c := range sr.Contributors {
		cs[i] = service.CreateContributorParam{
			FirstName: c.FirstName,
			LastName:  c.LastName,
			Email:     c.Email,
		}
	}

	return service.CreateSnippetParam{
		Title:        sr.Title,
		Code:         sr.Code,
		ProjectURL:   sr.ProjectURL,
		Git:          service.CreateGitParam{URL: sr.Git.URL},
		GitPath:      sr.GitPath,
		GitVersion:   sr.GitVersion,
		Language:     service.CreateLanguageParam{Name: sr.Language.Name},
		Tags:         ts,
		Contributors: cs,
	}
}

func toDeviceOAuth(r service.OAuthRedirect) DeviceOAuth {
	return DeviceOAuth{
		URL:       r.URL,
		UserCode:  r.UserCode,
		ExpiersIn: r.ExpiresIn,
		Interval:  r.Interval,
		Token:     r.Token,
	}
}

func toGithubPullRequest(r *http.Request) (GithubPullRequest, error) {
	defer r.Body.Close()

	var p GithubPullRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&p); err != nil {
		return GithubPullRequest{}, err
	}
	return p, nil
}

func toDeviceAuthResult(ar service.DeviceAuthResult) DeviceAuthResult {
	var session *Session
	if ar.Session != nil {
		session = &Session{
			User: User{
				ID:       strconv.FormatInt(ar.Session.User.ID, 10),
				Email:    ar.Session.User.Email,
				Username: ar.Session.User.Username,
			},
			TokenPair: TokenPair{
				AccessToken:  ar.Session.TokenPair.AccessToken,
				RefreshToken: ar.Session.TokenPair.RefreshToken,
			},
		}
	}
	return DeviceAuthResult{
		Status:  DeviceAuthStatus(ar.Status),
		Session: session,
	}
}
