package core

import (
	"context"
	"time"

	"github.com/osamikoyo/music-and-marks/services/mark/entity"
)

type Repository interface {
	CreateReview(ctx context.Context, review *entity.Review) error
	UpdateReview(ctx context.Context, id uint, update *entity.Review) error
	DeleteReview(ctx context.Context, id uint) error
	GetReviewsByReleaseID(ctx context.Context, releaseID string) ([]entity.Review, error)
	GetReviewByID(ctx context.Context, id uint) (*entity.Review, error)
	GetMarkByReleaseID(ctx context.Context, releaseID string) (*entity.Mark, error)
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

func NewCore(repo Repository, cache Cache, recounter Recounter, timeout time.Duration) *Core {
	return &Core{
		repo:      repo,
		cache:     cache,
		recounter: recounter,
		timeout:   timeout,
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

func (c *Core) IncLike(reviewID uint) error {
	ctx, cancel := c.context()
	defer cancel()

	review, err := c.repo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	review.Likes++

	return c.repo.UpdateReview(ctx, review.ID, review)
}

func (c *Core) DecLike(reviewID uint) error {
	ctx, cancel := c.context()
	defer cancel()

	review, err := c.repo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	review.Likes--

	return c.repo.UpdateReview(ctx, review.ID, review)
}
