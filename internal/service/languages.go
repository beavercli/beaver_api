package service

import (
	"context"

	"github.com/beavercli/beaver_api/internal/storage"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetLanguagesPage(ctx context.Context, p PageParam) (LanguageList, error) {
	var langs []storage.Language
	var total int64

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		langs, err = s.db.ListLanguages(ctx, storage.ListLanguagesParams{
			Offset: int32(p.Offset()),
			Limit:  int32(p.Limit()),
		})
		return err
	})
	g.Go(func() error {
		var err error
		total, err = s.db.CountLanguages(ctx)
		return err
	})
	if err := g.Wait(); err != nil {
		return LanguageList{}, err
	}
	return LanguageList{
		Total: int(total),
		Items: toServiceLanguage(langs),
	}, nil
}

func toServiceLanguage(ls []storage.Language) []Language {
	languages := make([]Language, len(ls))
	for i, l := range ls {
		languages[i] = Language{
			ID:   l.ID,
			Name: l.Name.String,
		}
	}
	return languages
}
