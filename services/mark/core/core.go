package core

import (
	"context"
	"time"

	"github.com/osamikoyo/music-and-marks/services/mark/entity"
)

type Repository interface {
	CreateReview(ctx context.Context, review *entity.Review) error
	DeleteReview(ctx context.Context, id uint) error
	GetReviewsByReleaseID(ctx context.Context, releaseID string) error
}

type Core struct {
	repo    Repository
	timeout time.Duration
}

func NewCore(repo Repository, timeout time.Duration) *Core {
	return &Core{
		repo:    repo,
		timeout: timeout,
	}
}

func (c *Core) CreateReview(releaseID, text, userID string, count int) error {
}
