package recounter

import (
	"context"
	"sync"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/mark/cache"
	"github.com/osamikoyo/music-and-marks/services/mark/core"
	"github.com/osamikoyo/music-and-marks/services/mark/entity"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const CountToRecount = 10

type Recounter struct {
	cache    *cache.Cache
	logger   *logger.Logger
	repo     core.Repository
	releases chan string

	mu     sync.RWMutex
	counts map[string]int
}

func NewRecounter(cache *cache.Cache, repo core.Repository, logger *logger.Logger) *Recounter {
	return &Recounter{
		cache:  cache,
		logger: logger,
		repo:   repo,
		counts: make(map[string]int),
	}
}

func (r *Recounter) recountMark(ctx context.Context, releaeID string) error {
	reviews, err := r.repo.GetReviewsByReleaseID(ctx, releaeID)
	if err != nil {
		return err
	}

	sum := 0

	for _, review := range reviews {
		sum += review.Count
	}

	count := float32(sum / len(reviews))

	mark := entity.NewMark(releaeID, count)

	if err := r.repo.UpdateMarkByReleaseID(ctx, releaeID, mark); err != nil {
		return err
	}

	return nil
}

func (r *Recounter) Start(appctx context.Context) {
	r.logger.Info("starting recounter...")

	eg, ctx := errgroup.WithContext(appctx)

	for {
		select {
		case <-appctx.Done():
			r.logger.Info("stopping recounter...")

			return
		case releaseID := <-r.releases:
			r.mu.Lock()
			r.counts[releaseID]++
			r.mu.Unlock()
		default:
			r.mu.Lock()
			for release, count := range r.counts {
				if count >= CountToRecount {
					eg.Go(func() error {
						return r.recountMark(ctx, release)
					})

					r.counts[release] = 0
				}
			}

			r.mu.Unlock()

			if err := eg.Wait(); err != nil {
				r.logger.Error("failed recount mark",
					zap.Error(err))
			}
		}
	}
}
