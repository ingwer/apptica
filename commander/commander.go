package commander

import (
	"apptica/fetcher"
	"apptica/model"
	"apptica/storage"

	"context"
)

type commander struct {
	fetcher fetcher.Fetcher
	storage storage.Storage
}

func (c commander) GetTopPositions(ctx context.Context, request Request) (map[model.Category]model.Position, error) {
	topPositions, err := c.storage.GetTopPositions(ctx, request.AppID, request.CountryCode, request.Date)
	if err != nil {
		return nil, err
	}

	if len(topPositions) == 0 {
		positions, err := c.fetcher.FetchTopPositions(ctx, request.AppID, request.CountryCode, request.Date, request.Date)
		if err != nil {
			return nil, err
		}

		for cat, catStat := range positions {
			err := c.storage.SaveTopPositions(ctx, request.AppID, request.CountryCode, cat, catStat)
			if err != nil {
				return nil, err
			}

			for date, position := range catStat {
				if date == request.Date {
					topPositions[cat] = position
				}
			}
		}
	}

	return topPositions, nil
}

type Request struct {
	CountryCode int
	AppID       int
	Date        model.Date
}

type Commander interface {
	GetTopPositions(ctx context.Context, request Request) (map[model.Category]model.Position, error)
}

type Option func(c *commander)

func New(options ...Option) Commander {
	c := &commander{}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithFetcher(f fetcher.Fetcher) Option {
	return func(c *commander) {
		c.fetcher = f
	}
}

func WithStorage(s storage.Storage) Option {
	return func(c *commander) {
		c.storage = s
	}
}
