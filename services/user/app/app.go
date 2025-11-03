package app

import (
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/user/config"
	"github.com/osamikoyo/music-and-marks/services/user/core"
	"github.com/osamikoyo/music-and-marks/services/user/repository"
	"github.com/osamikoyo/music-and-marks/services/user/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	logger *logger.Logger
	grpc   *grpc.Server
}

func SetupApp(config_path string) (*App, error) {
	logger.Init(logger.Config{
		AppName:   "user-service",
		AddCaller: false,
		LogFile:   "logs/user-service.log",
		LogLevel:  "debug",
	})

	logger := logger.Get()

	logger.Info("setuping app...")

	cfg, err := config.NewConfig(config_path, logger)
	if err != nil {
		logger.Error("failed load config",
			zap.String("path", config_path),
			zap.Error(err))

		return nil, fmt.Errorf("failed load config: %v", err)
	}

	logger.Info("setuping db...",
		zap.String("path", cfg.DatabasePath))

	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath))
	if err != nil {
		logger.Error("failed to connect to db",
			zap.String("database_path", cfg.DatabasePath),
			zap.Error(err))

		return nil, fmt.Errorf("failed setup database: %v", err)
	}

	repo := repository.NewRepository(db, logger)
	core := core.NewUserCore(repo, )
}
