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

func (r *Repository) Search(ctx context.Context, query string, pageSize, pageIndex int) ([]entity.SearchResult, error) {
	if query == "" {
		return nil, ErrNilInput
	}

	queryParam := query
	likeQuery := "%" + query + "%"
	offset := (pageIndex - 1) * pageSize

	var results []entity.SearchResult

	err := r.db.WithContext(ctx).Raw(`
        SELECT 
            id,
            mbid,
            title,
            artist_name,
            type,
            release_date,
            COALESCE(avg_rating, 0) AS avg_rating,
            COALESCE(review_count, 0) AS review_count,
            relevance
        FROM (
            -- Поиск по релизам (Release + ReleaseGroup)
            SELECT 
                rel.id::text,
                rel.mbid,
                rg.title,
                a.name AS artist_name,
                'release' AS type,
                rel.date AS release_date,
                AVG(rev.rating) AS avg_rating,
                COUNT(rev.id) AS review_count,
                GREATEST(
                    similarity(rg.title, ?),
                    similarity(a.name, ?)
                ) AS relevance
            FROM releases rel
            JOIN release_groups rg ON rel.release_group_id = rg.id
            JOIN artists a ON rg.artist_id = a.id
            LEFT JOIN reviews rev ON rev.target_id = rel.id AND rev.target_type = 'release'
            WHERE rg.title ILIKE ? 
               OR a.name ILIKE ?
               OR similarity(rg.title, ?) > 0.2
               OR similarity(a.name, ?) > 0.2
            GROUP BY rel.id, rel.mbid, rg.title, a.name, rel.date

            UNION ALL

            -- Поиск по артистам
            SELECT 
                a.id::text,
                a.mbid,
                a.name AS title,
                NULL::text AS artist_name,
                'artist' AS type,
                NULL::date AS release_date,
                NULL::float AS avg_rating,
                NULL::bigint AS review_count,
                similarity(a.name, ?) AS relevance
            FROM artists a
            WHERE a.name ILIKE ?
               OR similarity(a.name, ?) > 0.3
        ) combined
        WHERE relevance > 0.15
        ORDER BY relevance DESC, avg_rating DESC NULLS LAST
        LIMIT ? OFFSET ?
    `,
		queryParam, queryParam,
		likeQuery, likeQuery,
		queryParam, queryParam,
		queryParam,
		likeQuery,
		queryParam,
		pageSize, offset,
	).Scan(&results).Error
	if err != nil {
		r.logger.Error("failed to execute unified search", zap.Error(err))
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

func (r *Repository) GetReleaseByID(ctx context.Context, id uuid.UUID) (*entity.Release, error) {
	if len(id) == 0 {
		return nil, ErrNilInput
	}

	r.logger.Info("fetching release",
		zap.String("id", id.String()))

	var release entity.Release

	if err := r.db.WithContext(ctx).First(&release, id).Error; err != nil {
		r.logger.Error("failed fetch release",
			zap.String("id", id.String()),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	r.logger.Info("release successfully fetched",
		zap.Any("release", release))

	return &release, nil
}

func (r *Repository) ReadArtists(ctx context.Context, pageSize, pageIndex int) ([]entity.Artist, error) {
	r.logger.Info("fetching artists",
		zap.Int("pageSize", pageSize),
		zap.Int("pageIndex", pageIndex))

	var artists []entity.Artist

	offset := (pageIndex - 1) * pageSize

	res := r.db.Limit(pageSize).Offset(offset).Find(&artists)
	if res.RowsAffected == 0 {
		r.logger.Error("artists not found")

		return nil, ErrNotFound
	}
	if err := res.Error; err != nil {
		r.logger.Error("failed fetch artists",
			zap.Error(err))

		return nil, ErrInternal
	}

	r.logger.Info("artists have been successfully fetched")

	return artists, nil
}

func (r *Repository) ReadReleases(ctx context.Context, pageSize, pageIndex int) ([]entity.Release, error) {
	r.logger.Info("fetching releases",
		zap.Int("pageSize", pageSize),
		zap.Int("pageIndex", pageIndex))

	var releases []entity.Release

	offset := (pageIndex - 1) * pageSize

	res := r.db.Limit(pageSize).Offset(offset).Find(&releases)
	if res.RowsAffected == 0 {
		r.logger.Error("releases not found")

		return nil, ErrNotFound
	}
	if err := res.Error; err != nil {
		r.logger.Error("failed fetch release",
			zap.Error(err))

		return nil, ErrInternal
	}

	r.logger.Info("releases have been successfully fetched")

	return releases, nil
}
