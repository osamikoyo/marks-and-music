package core

import (
	"context"
	"time"

	"github.com/osamikoyo/music-and-marks/services/mark/entity"
)

type Repository interface {
	CreateReview(ctx context.Context, review *entity.Review) error
	DeleteReview(ctx context.Context, id uint) error
	GetReviewsByReleaseID(ctx context.Context, releaseID string) ([]entity.Review, error)
	GetMarkByReleaseID(ctx context.Context, releaseID string) (*entity.Mark, error)
	CreateMark(ctx context.Context, mark *entity.Mark) error
	UpdateMarkByReleaseID(ctx context.Context, releaseID string, update *entity.Mark) error
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

func (c *Core) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *Core) CreateReview(releaseID, text, userID string, count int) error {
	review := entity.NewReview(releaseID, text, userID, count)

	ctx, cancel := c.context()
	defer cancel()

	if err := c.repo.CreateReview(ctx, review); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetReviewsByReleaseID(releaeID string) ([]entity.Review, error) {
	ctx, cancel := c.context()
	defer cancel()

	reviews, err := c.repo.GetReviewsByReleaseID(ctx, releaeID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (c *Core) GetMarkByReleaeID(releaseID string) (*entity.Mark, error) {
	ctx, cancel := c.context()
	defer cancel()

	mark, err := c.repo.GetMarkByReleaseID(ctx, releaseID)
	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (c *Core) CreateMark(releaseID string, value float32) error {
	ctx, cancel := c.context()
	defer cancel()

	mark := entity.NewMark(releaseID, value)
	if err := c.repo.CreateMark(ctx, mark); err != nil {
		return err
	}

	return nil
}
