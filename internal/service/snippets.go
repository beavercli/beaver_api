package service

import (
	"context"
	"fmt"

	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetSnippet(ctx context.Context, id int64) (Snippet, error) {
	var tags []storage.GetTagsBySnippetIDRow
	var contributors []storage.GetContributorsBySnippetIDRow
	var snippet storage.GetSnippetByIDRow

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		snippet, err = s.db.GetSnippetByID(ctx, id)
		return err
	})
	g.Go(func() error {
		var err error
		tags, err = s.db.GetTagsBySnippetID(ctx, id)
		return err
	})
	g.Go(func() error {
		var err error
		contributors, err = s.db.GetContributorsBySnippetID(ctx, id)
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
		Tags:         convTags(tags),
		Contributors: convContributors(contributors),
	}, nil
}

type ListSnippetsParams struct {
	Page       int
	PageSize   int
	LanguageID *int64
	TagIDs     []int64
}

func (p ListSnippetsParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p ListSnippetsParams) Limit() int {
	return p.PageSize
}

func (s *Service) GetSnippetsPage(ctx context.Context, params ListSnippetsParams) (SnippetsList, error) {
	var snippetsCount int64
	var snippets []storage.ListSnippetsFilteredRow
	langID := pgtype.Int8{Valid: false}
	if params.LanguageID != nil {
		langID = pgtype.Int8{Int64: *params.LanguageID, Valid: true}
	}

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		snippetsCount, err = s.db.CountSnippetsFiltered(ctx, storage.CountSnippetsFilteredParams{
			LanguageID: langID,
			TagIds:     params.TagIDs,
		})
		if err != nil {
			fmt.Println("CountSnippetsFiltered error:", err) // Add this
		}
		return err
	})
	g.Go(func() error {
		var err error
		snippets, err = s.db.ListSnippetsFiltered(ctx, storage.ListSnippetsFilteredParams{
			LanguageID: langID,
			TagIds:     params.TagIDs,
			SqlLimit:   int32(params.Limit()),
			SqlOffset:  int32(params.Offset()),
		})
		if err != nil {
			fmt.Println("ListSnippetsFiltered error:", err) // Add this
		}
		return err
	})
	if err := g.Wait(); err != nil {
		fmt.Println("Error in waitgroup")
		return SnippetsList{}, err
	}

	snippetIDs := make([]int64, len(snippets))
	for _, s := range snippets {
		snippetIDs = append(snippetIDs, s.ID)
	}

	tags, err := s.db.GetTagsBySnippetIDs(ctx, snippetIDs)
	if err != nil {
		return SnippetsList{}, err
	}
	tagsBySnippet := mapTags(tags)

	snippetSummary := make([]SnippetSummary, len(snippets))
	for i, s := range snippets {
		snippetSummary[i] = SnippetSummary{
			ID:         s.ID,
			Title:      s.Title.String,
			ProjectURL: s.ProjectUrl.String,
			Language: &Language{
				ID:   s.LanguageID.Int64,
				Name: s.LanguageName.String,
			},
			Tags: tagsBySnippet[s.ID],
		}
	}

	return SnippetsList{
		Items: snippetSummary,
		Total: snippetsCount,
	}, nil
}

func mapTags(rows []storage.GetTagsBySnippetIDsRow) map[int64][]Tag {
	tagsBySnippet := make(map[int64][]Tag)

	for _, t := range rows {
		tagsBySnippet[t.SnippetID] = append(tagsBySnippet[t.SnippetID], Tag{ID: t.ID, Name: t.Name.String})
	}

	return tagsBySnippet
}

func convTags(rows []storage.GetTagsBySnippetIDRow) []Tag {
	tags := make([]Tag, len(rows))
	for i, r := range rows {
		tags[i] = Tag{ID: r.ID, Name: r.Name.String}
	}
	return tags
}

func convContributors(rows []storage.GetContributorsBySnippetIDRow) []Contributor {
	contributors := make([]Contributor, len(rows))
	for i, r := range rows {
		contributors[i] = Contributor{ID: r.ID, FirstName: r.FirstName.String, LastName: r.LastName.String}
	}
	return contributors
}
