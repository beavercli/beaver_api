package service

import (
	"context"

	"github.com/beavercli/beaver_api/internal/storage"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetContributorsPage(ctx context.Context, p PageParam) (ContributorList, error) {
	var contribs []storage.Contributor
	var total int64

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		contribs, err = s.db.ListContributors(ctx, storage.ListContributorsParams{
			Offset: int32(p.Offset()),
			Limit:  int32(p.Limit()),
		})
		return err
	})
	g.Go(func() error {
		var err error
		total, err = s.db.CountContributors(ctx)
		return err
	})
	if err := g.Wait(); err != nil {
		return ContributorList{}, err
	}
	return ContributorList{
		Total: int(total),
		Items: toServiceContributors(contribs),
	}, nil
}

func toServiceContributors(cs []storage.Contributor) []Contributor {
	contribs := make([]Contributor, len(cs))
	for i, c := range cs {
		contribs[i] = Contributor{
			ID:        c.ID,
			FirstName: c.FirstName.String,
			LastName:  c.LastName.String,
			Email:     c.Email.String,
		}
	}
	return contribs
}
