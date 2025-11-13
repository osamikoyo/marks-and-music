package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/loader"
	"github.com/osamikoyo/music-and-marks/services/music/repository"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultQueriesMaxCount = 5
	DefaultLoaderLimit     = 1
)

type Fetcher struct {
	logger  *logger.Logger
	loader  *loader.Loader
	repo    *repository.Repository
	queries chan string
	timeout time.Duration
}

func NewFetcher(loader *loader.Loader, repo *repository.Repository, logger *logger.Logger, timeout time.Duration) (*Fetcher, *FetcherClient) {
	queries := make(chan string, DefaultQueriesMaxCount)

	return &Fetcher{
		logger:  logger,
		repo:    repo,
		loader:  loader,
		timeout: timeout,
	}, newFetcherClient(queries)
}

func (f *Fetcher) Start(ctx context.Context) {
	f.logger.Info("starting async fetcher")

	for {
		select {
		case <-ctx.Done():
			f.logger.Info("fetcher stoping...")
			return
		case query := <-f.queries:
			f.logger.Info("new query for fetch",
				zap.String("query", query))

		}
	}
}

func (f *Fetcher) fetch(query string) error {
	f.logger.Info("fetching",
		zap.String("query", query))

	artists, err := f.loader.SearchArtists(query, 1)
	if err != nil {
		f.logger.Error("failed load artists",
			zap.String("query", query),
			zap.Error(err))

		return fmt.Errorf("failed load artists: %w", err)
	}

	artist := artists.Artists[0]

	albums, err := f.loader.SearchRelease(query, 1, 0)
	if err != nil {
		f.logger.Error("failed load releases",
			zap.String("query", query),
			zap.Error(err))

		return fmt.Errorf("failed load releaeses: %w", err)
	}

	album := albums.Releases[0]

	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return f.repo.CreateArtist(ctx, artist.ToEntity())
	})

	g.Go(func() error {
		
	})
}
