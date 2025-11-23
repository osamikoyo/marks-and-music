package core

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
)

var (
	ErrEmptyField = errors.New("empty field")
	ErrUIDFailed  = errors.New("failed parse uid")
	ErrPageIndex  = errors.New("page index less than zero")
)

type Repository interface {
	GetArtistByID(ctx context.Context, id uuid.UUID) (*entity.Artist, error)
	GetReleaseByID(ctx context.Context, id uuid.UUID) (*entity.Release, error)
	Search(ctx context.Context, query string, pageSize, pageIndex int) ([]entity.SearchResult, error)
	ReadArtists(ctx context.Context, pageSize, pageIndex int) ([]entity.Artist, error)
	ReadReleases(ctx context.Context, pageSize, pageIndex int) ([]entity.Release, error)
}

type Cache interface {
	SetArtist(key string, artist *entity.Artist) error
	SetRelease(key string, release *entity.Release) error
	GetArtist(key string) (*entity.Artist, error)
	GetRelease(key string) (*entity.Release, error)
}

type Fetcher interface {
	AsyncFetch(query string)
}

type MusicCore struct {
	cache   Cache
	fetcher Fetcher
	repo    Repository
	timeout time.Duration
}

func NewMusicCore(repo Repository, cache Cache, fetcher Fetcher, timeout time.Duration) *MusicCore {
	return &MusicCore{
		repo:    repo,
		cache:   cache,
		fetcher: fetcher,
		timeout: timeout,
	}
}

func (mc *MusicCore) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), mc.timeout)
}

func (mc *MusicCore) GetArtist(id string) (*entity.Artist, error) {
	if len(id) == 0 {
		return nil, ErrEmptyField
	}

	artist, err := mc.cache.GetArtist(id)
	if err == nil {
		return artist, nil
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUIDFailed
	}

	ctx, cancel := mc.context()
	defer cancel()

	artist, err = mc.repo.GetArtistByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	mc.cache.SetArtist(id, artist)

	return artist, nil
}

func (mc *MusicCore) GetRelease(id string) (*entity.Release, error) {
	if len(id) == 0 {
		return nil, ErrEmptyField
	}

	release, err := mc.cache.GetRelease(id)
	if err == nil {
		return release, nil
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUIDFailed
	}

	ctx, cancel := mc.context()
	defer cancel()

	release, err = mc.repo.GetReleaseByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	mc.cache.SetRelease(id, release)

	return release, nil
}

func (mc *MusicCore) Search(query string, pageIndex, pageSize int) ([]entity.SearchResult, error) {
	if len(query) == 0 {
		return nil, ErrEmptyField
	}

	ctx, cancel := mc.context()
	defer cancel()

	result, err := mc.repo.Search(ctx, query, pageSize, pageIndex)
	if err == nil {
		return result, nil
	}

	mc.fetcher.AsyncFetch(query)

	return nil, errors.New("not found results in db, fetching")
}

func (mc *MusicCore) ReadArtists(pageSize, pageIndex int) ([]entity.Artist, error) {
	if pageSize == 0 {
		return nil, nil
	}

	if pageIndex < 0 {
		return nil, ErrPageIndex
	}

	ctx, cancel := mc.context()
	defer cancel()

	artists, err := mc.repo.ReadArtists(ctx, pageSize, pageIndex)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (mc *MusicCore) ReadReleases(pageSize, pageIndex int) ([]entity.Release, error) {
	if pageSize == 0 {
		return nil, nil
	}

	if pageIndex < 0 {
		return nil, ErrPageIndex
	}

	ctx, cancel := mc.context()
	defer cancel()

	albums, err := mc.repo.ReadReleases(ctx, pageSize, pageIndex)
	if err != nil {
		return nil, err
	}

	return albums, nil
}
