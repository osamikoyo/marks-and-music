package server

import (
	"context"
	"errors"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/music/core"
	"github.com/osamikoyo/music-and-marks/services/music/metrics"
	"go.uber.org/zap"
)

var ErrEmptyRequest = errors.New("request is empty")

type Server struct {
	pb.UnimplementedMusicServiceServer
	core   *core.MusicCore
	logger *logger.Logger
}

func NewServer(core *core.MusicCore, logger *logger.Logger) *Server {
	return &Server{
		core:   core,
		logger: logger,
	}
}

func (s *Server) GetArtist(ctx context.Context, req *pb.GetArtistRequest) (*pb.GetArtistResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	then := time.Now()
	metrics.RequestTotal.WithLabelValues("GetArtist").Inc()

	artist, err := s.core.GetArtist(req.Id)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("GetArtist").Observe(time.Since(then).Seconds())

	return &pb.GetArtistResponse{
		Artist: artist.ToPB(),
	}, nil
}

func (s *Server) GetRelease(ctx context.Context, req *pb.GetReleaseRequest) (*pb.GetReleaseResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	then := time.Now()
	metrics.RequestTotal.WithLabelValues("GetRelease").Inc()

	release, err := s.core.GetRelease(req.Id)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("GetRelease").Observe(time.Since(then).Seconds())

	return &pb.GetReleaseResponse{
		Release: release.ToPB(),
	}, nil
}

func (s *Server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	then := time.Now()
	metrics.RequestTotal.WithLabelValues("Search").Inc()

	results, err := s.core.Search(req.Query, int(req.PageIndex), int(req.PageSize))
	if err != nil {
		s.logger.Error("failed search",
			zap.String("query", req.Query),
			zap.Error(err))

		return nil, err
	}

	pbresults := make([]*pb.SearchResult, len(results))
	for i, result := range results {
		pbresults[i] = result.ToPB()
	}

	metrics.RequestDuration.WithLabelValues("Search").Observe(time.Since(then).Seconds())

	return &pb.SearchResponse{
		Results: pbresults,
	}, nil
}

func (s *Server) ReadArtists(ctx context.Context, req *pb.ReadArtistsRequest) (*pb.ReadArtistsResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	then := time.Now()
	metrics.RequestTotal.WithLabelValues("ReadArtists").Inc()

	artists, err := s.core.ReadArtists(int(req.PageSize), int(req.PageIndex))
	if err != nil {
		s.logger.Error("failed read artists",
			zap.Int("page_size", int(req.PageSize)),
			zap.Int("page_index", int(req.PageIndex)),
			zap.Error(err))

		return nil, err
	}

	pbartists := make([]*pb.Artist, len(artists))
	for i, artist := range artists {
		pbartists[i] = artist.ToPB()
	}

	metrics.RequestDuration.WithLabelValues("ReadArtists").Observe(time.Since(then).Seconds())

	return &pb.ReadArtistsResponse{
		Artists: pbartists,
	}, nil
}

func (s *Server) ReadReleases(ctx context.Context, req *pb.ReadReleasesRequest) (*pb.ReadReleasesResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	then := time.Now()
	metrics.RequestTotal.WithLabelValues("ReadReleases")

	releases, err := s.core.ReadReleases(int(req.PageSize), int(req.PageIndex))
	if err != nil {
		s.logger.Error("failed read releases",
			zap.Int("page_size", int(req.PageSize)),
			zap.Int("page_index", int(req.PageIndex)),
			zap.Error(err))

		return nil, err
	}

	pbreleases := make([]*pb.Release, len(releases))
	for i, release := range releases {
		pbreleases[i] = release.ToPB()
	}

	metrics.RequestDuration.WithLabelValues("ReadReleases").Observe(time.Since(then).Seconds())

	return &pb.ReadReleasesResponse{
		Releases: pbreleases,
	}, nil
}
