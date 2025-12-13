package service

import (
	"context"
	"fmt"

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

func (s *Service) InjestSnippet(ctx context.Context, csp CreateSnippetParam) (Snippet, error) {
	txOptions := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}

	err := s.inTx(ctx, txOptions, func(db *storage.Queries) error {
		if err := s.uploadSnippetRelatedObjects(ctx, &csp); err != nil {
			return err
		}
		m, err := s.getMappingRelatedObjects(ctx, &csp)
		if err != nil {
			return err
		}
		// TODO
		fmt.Println(m)

		// get mappings to create the snippet

		return nil
	})

	if err != nil {
		return Snippet{}, err
	}
	return Snippet{}, nil
}

func (s *Service) uploadSnippetRelatedObjects(ctx context.Context, cs *CreateSnippetParam) error {
	g, upsertCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.db.UpsertLanguage(upsertCtx, pgtype.Text{String: cs.Language.Name, Valid: true}); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		tag_names := make([]string, len(cs.Tags))
		for i, t := range cs.Tags {
			tag_names[i] = t.Name
		}

		if err := s.db.BulkUpsertTags(upsertCtx, tag_names); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
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
		if err := s.db.BulkUpsertContributors(upsertCtx, contributorsParams); err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

type snippetMapping struct {
	tagsIDsByName     map[string]int64
	langIDsByName     map[string]int64
	contribsIDsByName map[string]int64
}

func (s *Service) getMappingRelatedObjects(ctx context.Context, cs *CreateSnippetParam) (snippetMapping, error) {
	m := snippetMapping{
		tagsIDsByName:     make(map[string]int64),
		langIDsByName:     make(map[string]int64),
		contribsIDsByName: make(map[string]int64),
	}
	g, ctxMapping := errgroup.WithContext(ctx)
	g.Go(func() error {
		tag_names := make([]string, len(cs.Tags))
		for i, t := range cs.Tags {
			tag_names[i] = t.Name
		}
		tags, err := s.db.GetTagIDsByNames(ctxMapping, tag_names)
		if err != nil {
			return err
		}
		for _, t := range tags {
			m.tagsIDsByName[t.Name.String] = t.ID
		}
		return nil
	})
	g.Go(func() error {
		emails := make([]string, len(cs.Contributors))
		for i, c := range cs.Contributors {
			emails[i] = c.Email
		}
		cs, err := s.db.GetContributorIDsByEmails(ctxMapping, emails)
		if err != nil {
			return err
		}
		for _, c := range cs {
			m.contribsIDsByName[c.Email.String] = c.ID
		}
		return nil
	})
	g.Go(func() error {
		id, err := s.db.GetLanguageIDByName(ctxMapping, pgtype.Text{String: cs.Language.Name, Valid: true})
		if err != nil {
			return err
		}
		m.langIDsByName[cs.Language.Name] = id
		return nil
	})
	if err := g.Wait(); err != nil {
		return snippetMapping{}, err
	}
	return m, nil
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
