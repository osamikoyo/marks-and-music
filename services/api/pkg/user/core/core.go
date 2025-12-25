package core

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/api/config"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/user/client"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/user/handler"
	"github.com/osamikoyo/music-and-marks/services/user/api/proto/gen/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type UserCore struct {
	handler *handler.Handler
}

func SetupUserCore(cfg *config.Config, logger *logger.Logger) (*UserCore, error) {
	conn, err := grpc.NewClient(cfg.UserServiceAddr)
	if err != nil {
		logger.Error("failed connect to user service",
			zap.String("addr", cfg.UserServiceAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed connect to user: %w", err)
	}

	c := pb.NewUserServiceClient(conn)
	client := client.NewUserClient(c, logger)
	handler := handler.NewHandler(client)

	return &UserCore{
		handler: handler,
	}, nil
}

func (u *UserCore) RegisterHandler(e *echo.Echo) {
}
