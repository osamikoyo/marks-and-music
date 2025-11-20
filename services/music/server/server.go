package server

import (
	"context"
	"errors"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/music/core"
)

var (
	ErrEmptyRequest = errors.New("request is empty")
)

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

	artist, err := s.core.GetArtist(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetArtistResponse{
		Artist: artist.ToPB(),
	}, nil
}

func (s *Server) GetRelease(ctx context.Context, req *pb.GetReleaseRequest) (*pb.GetReleaseResponse, error) {
	if req == nil{
		return nil, ErrEmptyRequest
	}

	release, err := s.core.GetRelease(req.Id)
	if err != nil{
		return nil, err
	}

	return &pb.GetReleaseResponse{
		Release: release.ToPB(),
	}, nil
}

func (s *Server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {

}

func (s *Server) ReadArtists(ctx context.Context, req *pb.ReadArtistsRequest) (*pb.ReadArtistsResponse, error) {

}

func (s *Server) ReadReleases(ctx context.Context, req *pb.ReadReleasesRequest) (*pb.ReadReleasesResponse, error) {

}
