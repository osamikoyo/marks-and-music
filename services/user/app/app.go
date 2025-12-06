package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/user/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/user/config"
	"github.com/osamikoyo/music-and-marks/services/user/core"
	"github.com/osamikoyo/music-and-marks/services/user/metrics"
	"github.com/osamikoyo/music-and-marks/services/user/repository"
	"github.com/osamikoyo/music-and-marks/services/user/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	if err = cfg.Validate(); err != nil {
		logger.Error("invalid config",
			zap.Error(err))

		return nil, fmt.Errorf("invalid config: %v", err)
	}

	db, err := setupDB(logger, cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(db, logger)
	core := core.NewUserCore(repo, cfg)
	server := server.NewUserServiceServer(core, logger)

	grpcSrv := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcSrv, server)

	logger.Info("server setup successfully")

	return &App{
		grpc:   grpcSrv,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func setupDB(logger *logger.Logger, cfg *config.Config) (*gorm.DB, error) {
	logger.Info("setuping db...",
		zap.String("path", cfg.DatabasePath))

	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath))
	if err != nil {
		logger.Error("failed to connect to db",
			zap.String("database_path", cfg.DatabasePath),
			zap.Error(err))

		return nil, fmt.Errorf("failed setup database: %v", err)
	}

	logger.Info("successfully setup db")

	return db, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("starting app...")

	metrics.InitMetrics()

	lis, err := net.Listen("tcp", a.cfg.Addr)
	if err != nil {
		a.logger.Error("failed listen",
			zap.String("addr", a.cfg.Addr),
			zap.Error(err))

		return fmt.Errorf("failed listen on %s: %v", a.cfg.Addr, err)
	}

	go func() {
		<-ctx.Done()

		a.grpc.GracefulStop()
	}()

	var wg sync.WaitGroup

	http.Handle("/metrics", promhttp.Handler())

	wg.Go(func() {
		if err = http.ListenAndServe(a.cfg.MetricsAddr, nil); err != nil {
			a.logger.Error("failed start metrics server",
				zap.String("addr", a.cfg.MetricsAddr),
				zap.Error(err))
		}
	})

	a.logger.Info("metrics server started",
		zap.String("addr", a.cfg.MetricsAddr))

	wg.Go(func() {
		if err = a.grpc.Serve(lis); err != nil {
			a.logger.Error("failed serve",
				zap.Error(err))
		}
	})

	a.logger.Info("grpc server started",
		zap.String("addr", a.cfg.Addr))

	wg.Wait()

	a.logger.Info("user-service started successfully")

	return nil
}
