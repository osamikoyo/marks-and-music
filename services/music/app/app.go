package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/music/cache"
	"github.com/osamikoyo/music-and-marks/services/music/config"
	"github.com/osamikoyo/music-and-marks/services/music/core"
	"github.com/osamikoyo/music-and-marks/services/music/fetcher"
	"github.com/osamikoyo/music-and-marks/services/music/loader"
	"github.com/osamikoyo/music-and-marks/services/music/repository"
	"github.com/osamikoyo/music-and-marks/services/music/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	fetcher    *fetcher.Fetcher
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

	logger.Info("config loaded successfully",
		zap.Any("config", cfg))

	dsn, err := cfg.Postgres.GetDSN()
	if err != nil {
		logger.Error("failed build dsn",
			zap.Any("db_config", cfg.Postgres),
			zap.Error(err))

		return nil, fmt.Errorf("build dsn error: %w", err)
	}

	logger.Info("dsn created",
		zap.String("dsn", dsn))

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		logger.Error("failed connect to database",
			zap.String("dsn", dsn),
			zap.Error(err))

		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	logger.Info("connected to database")

	repo := repository.NewRepository(db, logger)
	cache := cache.NewCache(cfg, logger)

	loader := loader.NewLoader(logger, cfg.SearchRequestTimeout)

	fetcher, fclient := fetcher.NewFetcher(loader, repo, logger, cfg.SearchRequestTimeout)
	core := core.NewMusicCore(repo, cache, fclient, cfg.RepositoryTimeout)

	server := server.NewServer(core, logger)
	grpcsrv := grpc.NewServer()
	pb.RegisterMusicServiceServer(grpcsrv, server)

	logger.Info("app setuped successfully")

	return &App{
		fetcher:    fetcher,
		grpcServer: grpcsrv,
		logger:     logger,
		cfg:        cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("starting app")

	var wg sync.WaitGroup

	wg.Go(func() {
		a.fetcher.Start(ctx)

		a.logger.Info("fetcher started")
	})

	lis, err := net.Listen("tcp", a.cfg.Addr)
	if err != nil {
		a.logger.Error("failed listen",
			zap.String("addr", a.cfg.Addr),
			zap.Error(err))

		return fmt.Errorf("failed listen: %w", err)
	}

	wg.Go(func() {
		if err = a.grpcServer.Serve(lis); err != nil {
			a.logger.Error("failed start grpc server",
				zap.Error(err))
		}
	})

	http.Handle("/metrics", promhttp.Handler())

	wg.Go(func() {
		if err = http.ListenAndServe(a.cfg.MetricsAddr, nil); err != nil {
			a.logger.Error("failed start metrics server",
				zap.String("addr", a.cfg.MetricsAddr),
				zap.Error(err))
		}
	})

	go func() {
		<-ctx.Done()

		a.grpcServer.GracefulStop()

		a.logger.Info("grpc server successfully stoped")
	}()

	wg.Wait()

	return nil
}
