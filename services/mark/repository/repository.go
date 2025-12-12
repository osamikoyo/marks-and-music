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
	ErrNotFound     = errors.New("not found")
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

func (r *Repository) DeleteReview(ctx context.Context, id uint) error {
	r.logger.Info("deleting review",
		zap.Uint("id", id))

	res := r.db.WithContext(ctx).Delete(&entity.Review{}, id)
	if err := res.Error; err != nil {
		r.logger.Error("failed delete review by release_id",
			zap.Uint("id", id),
			zap.Error(err))

		if res.RowsAffected == 0 {
			return ErrNotFound
		}

		return ErrInternal
	}

	r.logger.Info("review delete",
		zap.Uint("id", id))

	return nil
}

func (r *Repository) GetReviewsByReleaseID(ctx context.Context, releaseID string) ([]entity.Review, error) {
	r.logger.Info("fetching reviews by release id",
		zap.String("id", releaseID))

	var reviews []entity.Review
	res := r.db.WithContext(ctx).Where("release_id = ?", releaseID).Find(&reviews)

	if err := res.Error; err != nil {
		r.logger.Error("failed fetch release id",
			zap.String("release_id", releaseID),
			zap.Error(err))

		return nil, ErrInternal
	}

	r.logger.Info("reviews successfully fetched",
		zap.Int("len", len(reviews)))

	return reviews, nil
}

func (r *Repository) GetReviewByID(ctx context.Context, id uint) (*entity.Review, error) {
	r.logger.Info("fetching review",
		zap.Uint("id", id))

	var review entity.Review
	res := r.db.WithContext(ctx).First(&review, id)
	if err := res.Error; err != nil {
		r.logger.Error("failed fetch review",
			zap.Uint("id", id),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	r.logger.Info("review fetched successfully",
		zap.Any("review", review))

	return &review, nil
}

func (r *Repository) GetMarkByReleaseID(ctx context.Context, releaseID string) (*entity.Mark, error) {
	r.logger.Info("fetching mark",
		zap.String("release_id", releaseID))

	var mark entity.Mark
	res := r.db.WithContext(ctx).Where("release_id = ?", releaseID).First(&mark)
	if err := res.Error; err != nil {
		r.logger.Error("failed fetch mark",
			zap.String("release_id", releaseID),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	r.logger.Info("mark fetched",
		zap.Any("mark", mark))

	return &mark, nil
}

func (r *Repository) UpdateMarkByReleaseID(ctx context.Context, releaseID string, update *entity.Mark) error {
	r.logger.Info("updating mark",
		zap.String("release_id", releaseID),
		zap.Any("update", update))

	res := r.db.WithContext(ctx).Save(&update)

	if err := res.Error; err != nil {
		r.logger.Error("failed update mark",
			zap.String("release_id", releaseID),
			zap.Any("update", update),
			zap.Error(err))

		return ErrInternal
	}

	r.logger.Info("successfully updated mark",
		zap.String("release_id", releaseID),
		zap.Any("update", update))

	return nil
}
