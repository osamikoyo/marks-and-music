package core

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
)

const DefaultPageSize = 10

var (
	ErrEmptyField = errors.New("empty field")
	ErrUIDFailed  = errors.New("failed parse uid")
	ErrPageIndex  = errors.New("page index less than zero")
)

type Repository interface {
	GetArtistByID(ctx context.Context, id uuid.UUID) (*entity.Artist, error)
	GetReleaseByID(ctx context.Context, id uuid.UUID) (*entity.Release, error)
	Search(ctx context.Context, query string, page_size, page_index int) ([]entity.SearchResult, error)
	ReadArtists(ctx context.Context, page_size, page_index int) ([]entity.Artist, error)
	ReadReleases(ctx context.Context, page_size, page_index int) ([]entity.Release, error)
}

type Fetcher interface {
	AsyncFetch(query string)
}

type MusicCore struct {
	fetcher Fetcher
	repo    Repository
	timeout time.Duration
}

func (mc *MusicCore) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), mc.timeout)
}

func (mc *MusicCore) GetArtist(id string) (*entity.Artist, error) {
	if len(id) == 0 {
		return nil, ErrEmptyField
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUIDFailed
	}

	ctx, cancel := mc.context()
	defer cancel()
	return mc.repo.GetArtistByID(ctx, uid)
}

func (mc *MusicCore) GetRelease(id string) (*entity.Release, error) {
	if len(id) == 0 {
		return nil, ErrEmptyField
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUIDFailed
	}

	ctx, cancel := mc.context()
	defer cancel()
	return mc.repo.GetReleaseByID(ctx, uid)
}

func (mc *MusicCore) Search(query string, page_index int) ([]entity.AlbumSearchResult, error) {
	if len(query) == 0 {
		return nil, ErrEmptyField
	}

	ctx, cancel := mc.context()
	defer cancel()

	result, err := mc.repo.Search(ctx, query, DefaultPageSize, page_index)
	if err == nil {
		return result, nil
	}

	mc.fetcher.AsyncFetch(query)

	return nil, errors.New("not found results in db, fetching")
}

func (mc *MusicCore) ReadArtists(page_size, page_index int) ([]entity.Artist, error) {
	if page_size == 0 {
		return nil, nil
	}

	if page_index < 0 {
		return nil, ErrPageIndex
	}

	ctx, cancel := mc.context()
	defer cancel()

	artists, err := mc.repo.ReadArtists(ctx, page_size, page_index)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (mc *MusicCore) ReadReleases(page_size, page_index int) ([]entity.Release, error) {
	if page_size == 0 {
		return nil, nil
	}

	if page_index < 0 {
		return nil, ErrPageIndex
	}

	ctx, cancel := mc.context()
	defer cancel()

	albums, err := mc.repo.ReadReleases(ctx, page_size, page_index)
	if err != nil {
		return nil, err
	}

	return albums, nil
}
