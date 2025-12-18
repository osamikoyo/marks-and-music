package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
	"go.uber.org/zap"
)

var ErrNilInput = errors.New("empty field")

type MusicClient struct {
	cc     pb.MusicServiceClient
	logger *logger.Logger
}

func NewMusicClient(cc pb.MusicServiceClient, logger *logger.Logger) *MusicClient {
	return &MusicClient{
		cc:     cc,
		logger: logger,
	}
}

func (c *MusicClient) GetArtist(ctx context.Context, id string) (*entity.Artist, error) {
	if id == "" {
		return nil, ErrNilInput
	}

	resp, err := c.cc.GetArtist(ctx, &pb.GetArtistRequest{Id: id})
	if err != nil {
		c.logger.Error("failed to fetch artist",
			zap.String("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch artist: %w", err)
	}

	artist := &entity.Artist{
		ID:       resp.Artist.Id,
		Name:     resp.Artist.Name,
		SortName: *resp.Artist.SortName,
		Country:  *resp.Artist.Country,
		Type:     *resp.Artist.Type,
	}

	return artist, nil
}

func (c *MusicClient) GetRelease(ctx context.Context, id string) (*entity.Release, error) {
	if id == "" {
		return nil, ErrNilInput
	}

	resp, err := c.cc.GetRelease(ctx, &pb.GetReleaseRequest{Id: id})
	if err != nil {
		c.logger.Error("failed to fetch release",
			zap.String("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}

	release := &entity.Release{
		ID:              resp.Release.Id,
		MBID:            resp.Release.Mbid,
		Title:           resp.Release.Title,
		ReleaseGroupID:  resp.Release.ReleaseGroupId,
		Status:          *resp.Release.Status,
		Country:         *resp.Release.Country,
		Date:            resp.Release.Date,
		Format:          *resp.Release.Format,
		TrackCount:      int(resp.Release.TrackCount),
	}

	return release, nil
}

func (c *MusicClient) Search(ctx context.Context, query string, pageIndex, pageSize int32) ([]entity.SearchResult, error) {
	if query == "" {
		return nil, ErrNilInput
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageIndex < 0 {
		pageIndex = 0
	}

	resp, err := c.cc.Search(ctx, &pb.SearchRequest{
		Query:     query,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		c.logger.Error("failed search",
			zap.String("query", query),
			zap.Int32("page_index", pageIndex),
			zap.Int32("page_size", pageSize),
			zap.Error(err))
		return nil, fmt.Errorf("failed search: %w", err)
	}

	results := make([]entity.SearchResult, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = entity.SearchResult{
			ID:           r.Id,
			MBID:         r.Mbid,
			Title:        r.Title,
			ArtistName:   r.ArtistName,
			Type:         r.Type,
			ReleaseDate:  r.ReleaseDate,
			Relevance:    r.Relevance,
		}
	}

	return results, nil
}

func (c *MusicClient) ReadArtists(ctx context.Context, pageIndex, pageSize int32) ([]entity.Artist, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageIndex < 0 {
		pageIndex = 0
	}

	resp, err := c.cc.ReadArtists(ctx, &pb.ReadArtistsRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		c.logger.Error("failed read artists",
			zap.Int32("page_index", pageIndex),
			zap.Int32("page_size", pageSize),
			zap.Error(err))
		return nil, fmt.Errorf("failed read artists: %w", err)
	}

	artists := make([]entity.Artist, len(resp.Artists))
	for i, a := range resp.Artists {
		artists[i] = entity.Artist{
			ID:       a.Id,
			Name:     a.Name,
			SortName: *a.SortName,
			Country:  *a.Country,
			Type:     *a.Type,
		}
	}

	return artists, nil
}

func (c *MusicClient) ReadReleases(ctx context.Context, pageIndex, pageSize int32) ([]entity.Release, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageIndex < 0 {
		pageIndex = 0
	}

	resp, err := c.cc.ReadReleases(ctx, &pb.ReadReleasesRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		c.logger.Error("failed read releases",
			zap.Int32("page_index", pageIndex),
			zap.Int32("page_size", pageSize),
			zap.Error(err))
		return nil, fmt.Errorf("failed read releases: %w", err)
	}

	releases := make([]entity.Release, len(resp.Releases))
	for i, r := range resp.Releases {
		releases[i] = entity.Release{
			ID:              r.Id,
			MBID:            r.Mbid,
			Title:           r.Title,
			ReleaseGroupID:  r.ReleaseGroupId,
			Status:          *r.Status,
			Country:         *r.Country,
			Date:            r.Date,
			Format:          *r.Format,
			TrackCount:      int(r.TrackCount),
		}
	}

	return releases, nil
}