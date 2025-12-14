package service

import (
	"context"
	"fmt"
	"time"

	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetSnippet(ctx context.Context, id int64) (Snippet, error) {
	var tags []storage.GetTagsBySnippetIDRow
	var contributors []storage.GetContributorsBySnippetIDRow
	var snippet storage.GetSnippetByIDRow

	g, _ := errgroup.WithContext(ctx)

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
		GitRepoURL: snippet.GitRepoUrl.String,
		GitPath:    snippet.GitFilePath.String,
		GitVersion: snippet.GitVersion.String,
		Language: Language{
			ID:   snippet.LanguageID.Int64,
			Name: snippet.LanguageName.String,
		},
		Tags:         convTags(tags),
		Contributors: convContributors(contributors),
	}, nil
}

type ListSnippetsParams struct {
	PageParam

	LanguageID *int64
	TagIDs     []int64
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
	for i, s := range snippets {
		snippetIDs[i] = s.ID
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
			GitRepoURL: s.GitRepoUrl.String,
			GitPath:    s.GitFilePath.String,
			GitVersion: s.GitVersion.String,
			Language: Language{
				ID:   s.LanguageID.Int64,
				Name: s.LanguageName.String,
			},
			Tags: tagsBySnippet[s.ID],
		}
	}

	return SnippetsList{
		Items: snippetSummary,
		Total: int(snippetsCount),
	}, nil
}

type CreateContributorParam struct {
	FirstName string
	LastName  string
	Email     string
}
type CreateTagParam struct {
	Name string
}
type CreateLanguageParam struct {
	Name string
}
type CreateSnippetParam struct {
	Title        string
	Code         string
	ProjectURL   string
	GitRepoURL   string
	GitPath      string
	GitVersion   string
	Language     CreateLanguageParam
	Tags         []CreateTagParam
	Contributors []CreateContributorParam
}

func (s *Service) InjestSnippet(ctx context.Context, csp CreateSnippetParam) error {
	txOptions := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}

	err := s.inTx(ctx, txOptions, func(db *storage.Queries) error {
		r, err := uploadSnippetRelatedObjects(ctx, db, csp)
		if err != nil {
			return err
		}
		if err := updateOrCreateSnippet(ctx, db, csp, r); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

type snippetRefs struct {
	langID      int64
	tagsIDs     []int64
	contribsIDs []int64
}

func uploadSnippetRelatedObjects(ctx context.Context, tx *storage.Queries, cs CreateSnippetParam) (snippetRefs, error) {
	r := snippetRefs{
		tagsIDs:     make([]int64, len(cs.Tags)),
		contribsIDs: make([]int64, len(cs.Contributors)),
	}

	langID, err := tx.UpsertLanguage(ctx, pgtype.Text{String: cs.Language.Name, Valid: true})
	if err != nil {
		return snippetRefs{}, err
	}
	r.langID = langID

	tagNames := make([]string, len(cs.Tags))
	for i, t := range cs.Tags {
		tagNames[i] = t.Name
	}
	r.tagsIDs, err = tx.BulkUpsertTags(ctx, tagNames)
	if err != nil {
		return snippetRefs{}, err
	}

	contributorsParams := storage.BulkUpsertContributorsParams{
		FirstNames: make([]string, len(cs.Contributors)),
		LastNames:  make([]string, len(cs.Contributors)),
		Emails:     make([]string, len(cs.Contributors)),
	}
	for i, c := range cs.Contributors {
		contributorsParams.FirstNames[i] = c.FirstName
		contributorsParams.LastNames[i] = c.LastName
		contributorsParams.Emails[i] = c.Email
	}
	r.contribsIDs, err = tx.BulkUpsertContributors(ctx, contributorsParams)
	if err != nil {
		return snippetRefs{}, err
	}

	return r, nil
}

func updateOrCreateSnippet(ctx context.Context, tx *storage.Queries, cs CreateSnippetParam, r snippetRefs) error {
	snippetID, err := tx.UpsertSnippet(ctx, storage.UpsertSnippetParams{
		Title:       pgtype.Text{String: cs.Title, Valid: true},
		Code:        pgtype.Text{String: cs.Code, Valid: true},
		ProjectUrl:  pgtype.Text{String: cs.ProjectURL, Valid: true},
		GitRepoUrl:  pgtype.Text{String: cs.GitRepoURL, Valid: true},
		GitFilePath: pgtype.Text{String: cs.GitPath, Valid: true},
		GitVersion:  pgtype.Text{String: cs.GitVersion, Valid: true},
		LanguageID:  pgtype.Int8{Int64: r.langID, Valid: true},
		UserID:      pgtype.Int8{Valid: false},
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		return err
	}

	dr := storage.DeleteSnippetTagsExceptParams{
		SnippetID: snippetID,
		TagIds:    r.tagsIDs,
	}
	if err := tx.DeleteSnippetTagsExcept(ctx, dr); err != nil {
		return err
	}
	ur := storage.BulkLinkSnippetTagsParams{
		SnippetID: snippetID,
		TagIds:    r.tagsIDs,
	}
	if err := tx.BulkLinkSnippetTags(ctx, ur); err != nil {
		return err
	}

	drContrib := storage.DeleteSnippetContributorsExceptParams{
		SnippetID:      snippetID,
		ContributorIds: r.contribsIDs,
	}
	if err := tx.DeleteSnippetContributorsExcept(ctx, drContrib); err != nil {
		return err
	}
	urContrib := storage.BulkLinkSnippetContributorsParams{
		SnippetID:      snippetID,
		ContributorIds: r.contribsIDs,
	}
	if err := tx.BulkLinkSnippetContributors(ctx, urContrib); err != nil {
		return err
	}

	return nil
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
