package app

import (
	"fmt"
	"os"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/config"
	"github.com/osamikoyo/music-and-marks/services/music/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	grpcServer *grpc.Server
	logger     *logger.Logger
	cfg        *config.Config
}

func SetupApp() (*App, error) {
	logger.Init(logger.Config{
		AppName:   "music-service",
		LogFile:   "logs/music-service.log",
		LogLevel:  "debug",
		AddCaller: false,
	})

	logger := logger.Get()

	logger.Info("setupping app...")

	configpath := "music-service.yaml"
	for i, arg := range os.Args {
		if arg == "--config" {
			configpath = os.Args[i+1]
		}
	}

	cfg, err := config.NewConfig(configpath, logger)
	if err != nil {
		logger.Error("failed load config",
			zap.String("path", configpath),
			zap.Error(err))

		return nil, fmt.Errorf("load config error: %w", err)
	}

	dsn, err := cfg.Postgres.GetDSN()
	if err != nil {
		logger.Error("failed build dsn",
			zap.Any("db_config", cfg.Postgres),
			zap.Error(err))

		return nil, fmt.Errorf("build dsn error: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		logger.Error("failed connect to database",
			zap.String("dsn", dsn),
			zap.Error(err))

		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	repo := repository.NewRepository(db, logger)
}
