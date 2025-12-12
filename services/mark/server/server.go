package server

import (
	"context"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/mark/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/mark/core"
	"github.com/osamikoyo/music-and-marks/services/mark/metrics"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedMarkServiceServer
	logger *logger.Logger
	core   *core.Core
}

func NewServer(core *core.Core, logger *logger.Logger) *Server {
	return &Server{
		core:   core,
		logger: logger,
	}
}

func (s *Server) CreateReview(ctx context.Context, req *pb.Review) (*emptypb.Empty, error) {
	metrics.RequestTotal.WithLabelValues("CreateReview").Inc()
	then := time.Now()

	s.logger.Info("new create review request",
		zap.Any("req", req))

	if err := s.core.CreateReview(req.ReleaseId, req.Text, req.UserId, int(req.Count)); err != nil {
		return &emptypb.Empty{}, err
	}

	metrics.RequestDuration.WithLabelValues("CreateReview").Observe(time.Since(then).Seconds())

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*emptypb.Empty, error) {
	metrics.RequestTotal.WithLabelValues("DeleteReview").Inc()
	then := time.Now()

	s.logger.Info("new delete review request",
		zap.Any("req", req))

	if err := s.core.DeleteReview(uint(req.Id)); err != nil {
		return &emptypb.Empty{}, err
	}

	metrics.RequestDuration.WithLabelValues("DeleteReview").Observe(time.Since(then).Seconds())

	return &emptypb.Empty{}, nil
}

func (s *Server) GetMark(ctx context.Context, req *pb.GetMarkRequest) (*pb.Mark, error) {
	metrics.RequestTotal.WithLabelValues("GetMark").Inc()
	then := time.Now()

	s.logger.Info("new get mark request",
		zap.Any("req", req))

	mark, err := s.core.GetMarkByReleaeID(req.ReleaseId)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("GetMark").Observe(time.Since(then).Seconds())

	return mark.ToPB(), nil
}

func (s *Server) GetReviews(ctx context.Context, req *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	metrics.RequestTotal.WithLabelValues("GetReviews").Inc()
	then := time.Now()

	s.logger.Info("new get reviews request",
		zap.Any("req", req))

	reviews, err := s.core.GetReviewsByReleaseID(req.ReleaseId)
	if err != nil {
		return nil, err
	}

	pbreviews := make([]*pb.Review, len(reviews))
	for i, review := range reviews {
		pbreviews[i] = review.ToPB()
	}

	metrics.RequestDuration.WithLabelValues("GetReviews").Observe(time.Since(then).Seconds())

	return &pb.GetReviewsResponse{
		Reviews: pbreviews,
	}, nil
}
