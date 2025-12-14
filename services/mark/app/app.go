package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/mark/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/mark/cache"
	"github.com/osamikoyo/music-and-marks/services/mark/config"
	"github.com/osamikoyo/music-and-marks/services/mark/core"
	"github.com/osamikoyo/music-and-marks/services/mark/metrics"
	"github.com/osamikoyo/music-and-marks/services/mark/recounter"
	"github.com/osamikoyo/music-and-marks/services/mark/repository"
	"github.com/osamikoyo/music-and-marks/services/mark/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	grpc      *grpc.Server
	logger    *logger.Logger
	recounter *recounter.Recounter
	cfg       *config.Config
}

func SetupApp(configPath string) (*App, error) {
	logger.Init(logger.Config{
		AppName:   "mark-service",
		AddCaller: false,
		LogFile:   "logs/mark.log",
		LogLevel:  "debug",
	})

	logger := logger.Get()

	logger.Info("setupping app...",
		zap.String("config_path", configPath))

	cfg, err := config.NewConfig(configPath, logger)
	if err != nil {
		logger.Error("failed get config",
			zap.Error(err))

		return nil, fmt.Errorf("failed load config: %s: %w", configPath, err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBAddr))
	if err != nil {
		logger.Error("failed open db",
			zap.String("path", cfg.DBAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed open db: %s:%w", cfg.DBAddr, err)
	}

	logger.Info("connected to database",
		zap.String("db_addr", cfg.DBAddr))

	repo := repository.NewRepository(db, logger)

	cache := cache.NewCache(cfg, logger)

	recounter, client := recounter.NewRecounter(cache, repo, logger)

	core := core.NewCore(repo, cache, client, cfg.RepoTimeout)
	server := server.NewServer(core, logger)
	grpcsrv := grpc.NewServer()
	pb.RegisterMarkServiceServer(grpcsrv, server)

	metrics.InitMetrics()

	return &App{
		grpc:      grpcsrv,
		logger:    logger,
		recounter: recounter,
	}, nil
}

func (a *App) Run(appctx context.Context) error {
	a.logger.Info("starting app...")

	eg, ctx := errgroup.WithContext(appctx)

	lis, err := net.Listen("tcp", a.cfg.Addr)
	if err != nil {
		a.logger.Error("failed listen",
			zap.String("addr", a.cfg.Addr),
			zap.Error(err))

		return fmt.Errorf("failed listen: %w", err)
	}

	eg.Go(func() error {
		a.recounter.Start(ctx)

		return nil
	})

	http.Handle("/metrics", promhttp.Handler())

	eg.Go(func() error {
		if err := http.ListenAndServe(a.cfg.MetricsAddr, nil); err != nil {
			a.logger.Error("failed start metrics server",
				zap.String("addr", a.cfg.MetricsAddr),
				zap.Error(err))

			return fmt.Errorf("failed start metrics server: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := a.grpc.Serve(lis); err != nil {
			a.logger.Error("failed start grpc server",
				zap.String("addr", a.cfg.Addr),
				zap.Error(err))

			return fmt.Errorf("failed start grpc server: %w", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	a.logger.Info("mark service started successfully",
		zap.String("metrics_addr", a.cfg.MetricsAddr),
		zap.String("grpc_addr", a.cfg.Addr))

	return nil
}
