package service

import (
	"context"

	"github.com/beavercli/beaver_api/internal/storage"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetRandomSnippet(ctx context.Context) (Snippet, error) {
	var tags []storage.GetTagsBySnippetIDRow
	var contributors []storage.GetContributorsBySnippetIDRow
	var snippet storage.GetSnippetByIDRow

	snippetID := s.ids[0]

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		snippet, err = s.db.GetSnippetByID(ctx, snippetID)
		return err
	})
	g.Go(func() error {
		var err error
		tags, err = s.db.GetTagsBySnippetID(ctx, snippetID)
		return err
	})
	g.Go(func() error {
		var err error
		contributors, err = s.db.GetContributorsBySnippetID(ctx, snippetID)
		return err
	})

	if err := g.Wait(); err != nil {
		return Snippet{}, err
	}

	return Snippet{
		ID:         snippet.ID,
		Title:      snippet.Title.String,
		Code:       snippet.Code.String,
		ProjectURL: snippet.ProjectUrl.String,
		Language: &Language{
			ID:   snippet.LanguageID.Int64,
			Name: snippet.LanguageName.String,
		},
		Tags:         mapTags(tags),
		Contributors: mapContributors(contributors),
	}, nil
}

func mapTags(rows []storage.GetTagsBySnippetIDRow) []Tag {
	tags := make([]Tag, len(rows))
	for i, r := range rows {
		tags[i] = Tag{ID: r.ID, Name: r.Name.String}
	}
	return tags
}

func mapContributors(rows []storage.GetContributorsBySnippetIDRow) []Contributor {
	contributors := make([]Contributor, len(rows))
	for i, r := range rows {
		contributors[i] = Contributor{ID: r.ID, FirstName: r.FirstName.String, LastName: r.LastName.String}
	}
	return contributors
}
