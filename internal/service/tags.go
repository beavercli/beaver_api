package service

import (
	"context"

	"github.com/beavercli/beaver_api/internal/storage"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetTagsPage(ctx context.Context, p PageParam) (TagList, error) {
	var tags []storage.Tag
	var total int64

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		tags, err = s.db.ListTags(ctx, storage.ListTagsParams{
			Offset: int32(p.Offset()),
			Limit:  int32(p.Limit()),
		})
		return err
	})
	g.Go(func() error {
		var err error
		total, err = s.db.CountTags(ctx)
		return err
	})
	if err := g.Wait(); err != nil {
		return TagList{}, err
	}
	return TagList{
		Total: int(total),
		Items: toServiceTag(tags),
	}, nil
}

func toServiceTag(ts []storage.Tag) []Tag {
	tags := make([]Tag, len(ts))
	for i, t := range ts {
		tags[i] = Tag{
			ID:   t.ID,
			Name: t.Name.String,
		}
	}
	return tags
}
