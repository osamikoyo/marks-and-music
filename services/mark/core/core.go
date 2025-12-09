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
	GetReviewByID(ctx context.Context, id uint) (*entity.Review, error)
	GetMarkByReleaseID(ctx context.Context, releaseID string) (*entity.Mark, error)
	CreateMark(ctx context.Context, mark *entity.Mark) error
	UpdateMarkByReleaseID(ctx context.Context, releaseID string, update *entity.Mark) error
}

type Cache interface {
	Set(key string, value interface{})
	GetReviews(key string) ([]entity.Review, error)
}

type Recounter interface {
	TryRecount(releaseID string)
}

type Core struct {
	repo      Repository
	cache     Cache
	recounter Recounter
	timeout   time.Duration
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

	reviews, err := c.cache.GetReviews(releaseID)
	if err != nil {
		return err
	}

	newrevs := make([]entity.Review, len(reviews)+1)
	newrevs[len(reviews)+1] = *review

	c.cache.Set(releaseID, newrevs)

	c.recounter.TryRecount(releaseID)

	return nil
}

func (c *Core) GetReviewsByReleaseID(releaseID string) ([]entity.Review, error) {
	ctx, cancel := c.context()
	defer cancel()

	reviews, err := c.cache.GetReviews(releaseID)
	if err == nil {
		return reviews, nil
	}

	reviews, err = c.repo.GetReviewsByReleaseID(ctx, releaseID)
	if err != nil {
		return nil, err
	}

	c.cache.Set(releaseID, reviews)

	return reviews, nil
}

func (c *Core) DeleteReview(id uint) error {
	ctx, cancel := c.context()
	defer cancel()

	review, err := c.repo.GetReviewByID(ctx, id)
	if err != nil {
		return err
	}

	if err := c.repo.DeleteReview(ctx, id); err != nil {
		return err
	}

	c.recounter.TryRecount(review.ReleaseID)

	return nil
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
