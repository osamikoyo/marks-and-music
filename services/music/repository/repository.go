package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	logger *logger.Logger
	db     *gorm.DB
}

var (
	ErrInternal     = errors.New("internal error")
	ErrAlreadyExist = errors.New("user already exist")
	ErrNotFound     = errors.New("user not found")
	ErrNilInput     = errors.New("empty releasegroup")
)

func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
func (r *Repository) CreateReleaseGroup(ctx context.Context, rg *entity.ReleaseGroup) error {
	if rg == nil {
		return ErrNilInput
	}

	r.logger.Info("creating release group",
		zap.Any("rg", rg))

	if err := r.db.WithContext(ctx).Create(rg).Error; err != nil {
		r.logger.Error("failed create release group",
			zap.Error(err))

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExist
		}

		return ErrInternal
	}

	r.logger.Info("release group created successfully")

	return nil
}

func (r *Repository) CreateArtist(ctx context.Context, artist *entity.Artist) error {
	if artist == nil {
		return ErrNilInput
	}

	r.logger.Info("creating artist",
		zap.Any("artist", artist))

	if err := r.db.WithContext(ctx).Create(artist).Error; err != nil {
		r.logger.Error("failed create artist",
			zap.Error(err))

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExist
		}

		return ErrInternal
	}

	r.logger.Info("artist created successfully")

	return ErrNilInput
}

func (r *Repository) Search(ctx context.Context, query string, page_size, page_index int) ([]entity.AlbumSearchResult, error) {
	db := r.db.WithContext(ctx).Model(&entity.ReleaseGroup{}).
		Select(`
        rg.id,
        rg.mbid,
        rg.title,
        a.name AS artist_name,
        rg.first_release_date,
        COALESCE(AVG(r.rating), 0) AS avg_rating,
        COUNT(r.id) AS review_count,
        similarity(rg.title, ?) AS relevance
    `, query).
		Joins("JOIN artists a ON rg.artist_id = a.id").
		Joins("LEFT JOIN reviews r ON r.target_id = rg.id AND r.target_type = 'album'")

	db = db.Where(`
    rg.title ILIKE ? OR 
    similarity(rg.title, ?) > 0.3
`, "%"+query+"%", query)

	db = db.Group("rg.id").Group("a.name")

	db = db.Having("similarity(rg.title, ?) > 0.2", query)

	db = db.Order("relevance DESC").
		Order("avg_rating DESC")

	offset := (page_index - 1) * page_size

	db = db.Limit(page_size).Offset(offset)

	var results []entity.AlbumSearchResult
	if err := db.Scan(&results).Error; err != nil {
		r.logger.Error("failed scan search result",
			zap.Error(err))

		return nil, ErrInternal
	}

	return results, nil
}

func (r *Repository) CreateRelease(ctx context.Context, release *entity.Release) error {
	if release == nil {
		return ErrNilInput
	}

	r.logger.Info("creating release",
		zap.Any("release", release))

	if err := r.db.WithContext(ctx).Create(release).Error; err != nil {
		r.logger.Error("failed create artist",
			zap.Error(err))

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExist
		}

		return ErrInternal
	}

	r.logger.Info("release created successfully")

	return ErrNilInput
}

func (r *Repository) GetArtistByID(ctx context.Context, id uuid.UUID) (*entity.Artist, error) {
	if len(id) == 0 {
		return nil, ErrNilInput
	}

	r.logger.Info("fetching artist",
		zap.String("id", id.String()))

	var artist entity.Artist

	if err := r.db.WithContext(ctx).First(&artist, id).Error; err != nil {
		r.logger.Error("failed fetch artist",
			zap.String("id", id.String()),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	r.logger.Info("artist successfully fetched",
		zap.Any("artist", artist))
	return &artist, nil
}
