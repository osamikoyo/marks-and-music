package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/loader"
	"github.com/osamikoyo/music-and-marks/services/music/repository"
	"go.uber.org/zap"
)

const DefaultLoadLimit = 3

var (
	ErrPayloadType = errors.New("wrong payload type")
)

type Job struct {
	Type    string
	Payload any
}

type Worker struct {
	logger *logger.Logger
	repo   *repository.Repository
	loader *loader.Loader
	jobs   chan Job
}

func NewWorker(logger *logger.Logger, repo *repository.Repository, jobs chan Job) *Worker {
	return &Worker{
		logger: logger,
		repo:   repo,
		jobs:   jobs,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopped")

			return
		case job := <-w.jobs:

		}
	}
}

func (w *Worker) routeJob(job *Job) error {
	w.logger.Info("routing job...",
		zap.Any("job", job))

	switch job.Type {
	case "fetch":
		query, ok := job.Payload.(string)
		if !ok {
			return ErrPayloadType
		}

		artistSearchResult, err := w.loader.SearchArtists(query, DefaultLoadLimit)
		if err != nil {
			w.logger.Error("failed search artist",
				zap.String("query", query),
				zap.Error(err))

			return fmt.Errorf("failed search artist: %w", err)
		}

		releaseSearchResult, err := w.loader.SearchRelease(query, DefaultLoadLimit, 0)
		if err != nil {
			w.logger.Error("failed search result",
				zap.String("query", query),
				zap.Error(err))

			return fmt.Errorf("failed search release: %w", err)
		}
	}
}
