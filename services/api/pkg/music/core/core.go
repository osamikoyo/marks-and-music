package core

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/api/config"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/music/client"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/music/handler"
	"github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MusicCore struct {
	handler *handler.Handler
}

func SetupMusicCore(cfg *config.Config, logger *logger.Logger) (*MusicCore, error) {
	conn, err := grpc.NewClient(cfg.MarkServiceAddr)
	if err != nil {
		logger.Error("failed connect to mark service",
			zap.String("addr", cfg.MarkServiceAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed connect to mark service: %w", err)
	}

	cc := pb.NewMusicServiceClient(conn)
	client := client.NewMusicClient(cc, logger)

	handler := handler.NewHandler(client)

	return &MusicCore{
		handler: handler,
	}, nil
}

func (m *MusicCore) RegisterHandler(e *echo.Echo) {
	e.GET("/search", m.handler.Search)
	e.GET("/release/:id", m.handler.GetRelease)
	e.GET("/artist/:id", m.handler.GetArtist)
	e.GET("/releases", m.handler.ReadReleases)
	e.GET("/artists", m.handler.ReadArtists)
}

