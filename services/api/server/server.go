package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/api/config"
	markcore "github.com/osamikoyo/music-and-marks/services/api/pkg/mark/core"
	musiccore "github.com/osamikoyo/music-and-marks/services/api/pkg/music/core"
	usercore "github.com/osamikoyo/music-and-marks/services/api/pkg/user/core"
	"go.uber.org/zap"
)

type Core interface {
	RegisterHandler(e *echo.Echo)
}

type Server struct {
	e      *echo.Echo
	logger *logger.Logger
	cfg    *config.Config
}

func SetupApiServer(config_path string) (*Server, error) {
	logger.Init(logger.Config{
		AppName:   "api",
		AddCaller: false,
		LogFile:   "logs/api.log",
		LogLevel:  "debug",
	})

	logger := logger.Get()

	logger.Info("setupping api server")

	cfg, err := config.NewConfig(config_path, logger)
	if err != nil {
		logger.Error("faield get config",
			zap.String("path", config_path),
			zap.Error(err))

		return nil, fmt.Errorf("failed load config: %w", err)
	}

	cores, err := setupCores(cfg, logger)
	if err != nil {
		logger.Error("failed setup cores",
			zap.Error(err))

		return nil, fmt.Errorf("failed setup cores: %w", err)
	}

	e := echo.New()

	for _, core := range cores {
		core.RegisterHandler(e)
	}

	return &Server{
		e:      e,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting server",
		zap.String("addr", s.cfg.Addr))

	if err := s.e.Start(s.cfg.Addr); err != nil {
		s.logger.Error("failed start server",
			zap.Error(err))

		return fmt.Errorf("failed started server: %w", err)
	}

	return nil
}

func (s *Server) Close(ctx context.Context) error {
	s.logger.Info("closing server")

	return s.e.Shutdown(ctx)
}

func setupCores(cfg *config.Config, logger *logger.Logger) ([]Core, error) {
	logger.Info("setup cores")

	mark, err := markcore.SetupMarkCore(cfg, logger)
	if err != nil {
		logger.Error("failed setup mark core",
			zap.Error(err))

		return nil, fmt.Errorf("failed setup mark core: %w", err)
	}

	music, err := musiccore.SetupMusicCore(cfg, logger)
	if err != nil {
		logger.Error("failed setup music core",
			zap.Error(err))

		return nil, fmt.Errorf("failed setup music core: %w", err)
	}

	user, err := usercore.SetupUserCore(cfg, logger)
	if err != nil {
		logger.Error("failed setup user core",
			zap.Error(err))

		return nil, fmt.Errorf("failed setup user core: %w", err)
	}

	cores := make([]Core, 3)
	cores[1] = mark
	cores[2] = music
	cores[3] = user

	return cores, nil
}
