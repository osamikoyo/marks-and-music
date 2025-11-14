package core

import (
	"context"
	"errors"
	"time"

	"github.com/osamikoyo/music-and-marks/services/music/entity"
)

const DefaultPageSize = 10

var(
	ErrEmptyField = errors.New("empty field")
)

type Repository interface {
	GetArtistByID(ctx context.Context, id string) (*entity.Artist, error)
	GetReleaseByID(ctx context.Context) (*entity.Release, error)
	Search(ctx context.Context, query string, page_size, page_index int) ([]entity.AlbumSearchResult, error)
}

type Fetcher interface{
	AsyncFetch(query string)
}

type MusicCore struct{
	fetcher Fetcher
	repo Repository
	timeout time.Duration
}

func (mc *MusicCore) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), mc.timeout)
}

func (mc *MusicCore) GetArtist(id string) (*entity.Artist, error) {
	if len(id) == 0 {
		return nil, ErrEmptyField
	}

	ctx, cancel := mc.context()
	defer cancel()
	return mc.repo.GetArtistByID(ctx, id)
}

func (mc *MusicCore) Search(query string, page_index int) ([]entity.AlbumSearchResult, error) {
	if len(query) == 0 {
		return nil, ErrEmptyField
	}

	ctx, cancel := mc.context()
	defer cancel()

	result, err := mc.repo.Search(ctx, query, DefaultPageSize, page_index)
	if err == nil{
		return result, nil
	}

	mc.fetcher.AsyncFetch(query)

	return nil, errors.New("not found results in db, fetching")
}
