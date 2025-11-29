package repository

import (
	"context"
	"errors"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/mark/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrEmptyFields  = errors.New("empty fields")
	ErrInternal     = errors.New("internal error")
	ErrAlreadyExist = errors.New("already exist")
)

type Repository struct {
	logger *logger.Logger
	db     *gorm.DB
}

func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateReview(ctx context.Context, review *entity.Review) error {
	if review == nil {
		return ErrEmptyFields
	}

	r.logger.Info("creating review",
		zap.Any("review", review))

	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		r.logger.Error("failed create reciew",
			zap.Any("review", review),
			zap.Error(err))

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExist
		}

		return ErrInternal
	}

	r.logger.Info("review created successfully",
		zap.Any("review", review))

	return nil
}

func (r *Repository) DeleteReviewByReleaseID(ctx context.Context, releaseID string) error {
	if len(releaseID) == 0 {
		return ErrEmptyFields
	}

	r.logger.Info("deleting review",
		zap.String("release_id", releaseID))

	res := r.db.WithContext(ctx).Where("release_id = ?", releaseID).Delete(&entity.Review{})
}
