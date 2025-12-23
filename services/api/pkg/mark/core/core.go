package core

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/api/config"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/mark/client"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/mark/handler"
	"github.com/osamikoyo/music-and-marks/services/mark/api/proto/gen/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MarkCore struct {
	handler *handler.Handler
}

func SetupMarkCore(cfg *config.Config, logger *logger.Logger) (*MarkCore, error) {
	conn, err := grpc.NewClient(cfg.MarkServiceAddr)
	if err != nil {
		logger.Error("failed connect to mark service",
			zap.String("addr", cfg.MarkServiceAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed connect to mark service: %w", err)
	}

	cc := pb.NewMarkServiceClient(conn)
	client := client.NewMarkClient(cc, logger)

	handler := handler.NewHandler(client)

	return &MarkCore{
		handler: handler,
	}, nil
}

func (m *MarkCore) RegisterHandler(e *echo.Echo) {
	e.GET("/reviews/:releaseid", m.handler.GetReviews)
	e.GET("/mark/:releaseid", m.handler.GetMark)

	e.POST("/review/create", m.handler.CreateReview)

	e.DELETE("/review/delete", m.handler.DeleteReview)
}
