package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/mark/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/mark/entity"
	"go.uber.org/zap"
)

var ErrNilInput = errors.New("empty field")

type MarkClient struct {
	cc     pb.MarkServiceClient
	logger *logger.Logger
}

func NewMarkClient(cc pb.MarkServiceClient, logger *logger.Logger) *MarkClient {
	return &MarkClient{
		cc:     cc,
		logger: logger,
	}
}

func (u *MarkClient) GetReviews(ctx context.Context, releaseID string) ([]entity.Review, error) {
	if releaseID == "" {
		return nil, ErrNilInput
	}

	resp, err := u.cc.GetReviews(ctx, &pb.GetReviewsRequest{ReleaseId: releaseID})
	if err != nil {
		u.logger.Error("failed fetch reviews",
			zap.String("release_id", releaseID),
			zap.Error(err))

		return nil, fmt.Errorf("failed fetch reviews: %w", err)
	}

	reviews := make([]entity.Review, len(resp.Reviews))

	for i, review := range resp.Reviews {
		reviews[i] = entity.Review{
			ID:        uint(review.Id),
			UserID:    review.UserId,
			Text:      review.Text,
			Count:     int(review.Count),
			ReleaseID: review.ReleaseId,
		}
	}

	return reviews, nil
}

func (u *MarkClient) DeleteReview(ctx context.Context, id uint) error {
	_, err := u.cc.DeleteReview(ctx, &pb.DeleteReviewRequest{
		Id: uint64(id),
	})

	if err != nil {
		u.logger.Error("failed delete review",
			zap.Uint("id", id),
			zap.Error(err))

		return fmt.Errorf("failed delete review: %w", err)
	}

	return nil
}

func (u *MarkClient) GetMark(ctx context.Context, releaseID string) (*entity.Mark, error) {
	if releaseID == "" {
		return nil, ErrNilInput
	}

	pbmark, err := u.cc.GetMark(ctx, &pb.GetMarkRequest{ReleaseId: releaseID})
	if err != nil {
		u.logger.Error("failed fetch mark",
			zap.String("release_id", releaseID),
			zap.Error(err))

		return nil, fmt.Errorf("failed fetch mark: %w", err)
	}

	mark := entity.NewMark(pbmark.ReleaseId, pbmark.Value)

	return mark, nil
}

func (u *MarkClient) CreateReview(ctx context.Context, review *entity.Review) error {
	if review == nil {
		return ErrNilInput
	}

	_, err := u.cc.CreateReview(ctx, review.ToPB())
	if err != nil {
		u.logger.Error("failed create review",
			zap.Any("review", review),
			zap.Error(err))

		return fmt.Errorf("failed create review: %w", err)
	}

	return nil
}
