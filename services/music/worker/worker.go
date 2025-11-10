package worker

import (
	"context"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/repository"
	"go.uber.org/zap"
)

type Job struct {
	Type    string
	Payload any
}

type Worker struct {
	logger *logger.Logger
	repo   *repository.Repository
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
	case "fetch_artist":
		
		
	}
}
