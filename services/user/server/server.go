package server

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/user/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/user/core"
	"github.com/osamikoyo/music-and-marks/services/user/metrics"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrEmptyReq    = errors.New("request is empty")
	ErrInvalidUUID = errors.New("invalid uuid in request")
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	logger *logger.Logger
	core   *core.UserCore
}

func NewUserServiceServer(core *core.UserCore, logger *logger.Logger) *UserServiceServer {
	return &UserServiceServer{
		logger: logger,
		core:   core,
	}
}

func (uss *UserServiceServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	then := time.Now
	metrics.RequestTotal.WithLabelValues("ChangePassword").Inc()

	if req == nil {
		uss.logger.Error("empty request")

		return &emptypb.Empty{}, ErrEmptyReq
	}

	uss.logger.Info("new change password request",
		zap.String("old", req.CurrentPassword),
		zap.String("uid", req.Id),
		zap.String("new", req.NewPassword))

	uid, err := uuid.Parse(req.Id)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.Id))

		return &emptypb.Empty{}, ErrInvalidUUID
	}

	if err = uss.core.ChangePassword(uid, req.CurrentPassword, req.NewPassword); err != nil {
		return &emptypb.Empty{}, err
	}

	metrics.RequestDuration.WithLabelValues("ChangePassword").Observe(time.Since(then()).Seconds())

	return &emptypb.Empty{}, nil
}

func (uss *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	then := time.Now()
	metrics.RequestTotal.WithLabelValues("DeleteUser").Inc()
	if req == nil {
		uss.logger.Error("empty request")

		return &emptypb.Empty{}, ErrEmptyReq
	}

	uss.logger.Info("new delete user request",
		zap.String("id", req.Id))

	uid, err := uuid.Parse(req.Id)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.Id))

		return &emptypb.Empty{}, ErrInvalidUUID
	}

	if err := uss.core.DeleteUser(uid); err != nil {
		return &emptypb.Empty{}, err
	}

	metrics.RequestDuration.WithLabelValues("DeleteUser").Observe(time.Since(then).Seconds())

	return &emptypb.Empty{}, nil
}

func (uss *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	then := time.Now()
	metrics.RequestTotal.WithLabelValues("GetUser").Inc()

	if req == nil {
		uss.logger.Error("empty request")

		return nil, ErrEmptyReq
	}

	uss.logger.Info("new get user request",
		zap.String("id", req.Id))

	uid, err := uuid.Parse(req.Id)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.Id))

		return nil, ErrInvalidUUID
	}

	user, err := uss.core.GetUserByID(uid)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("GetUser").Observe(float64(time.Since(then).Seconds()))

	return user.ToProto(), nil
}

func (uss *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.TokenPair, error) {
	then := time.Now()
	metrics.RequestTotal.WithLabelValues("Login").Inc()

	if req == nil {
		uss.logger.Error("empty request")

		return nil, ErrEmptyReq
	}

	uss.logger.Info("new login request",
		zap.String("email", req.Email),
		zap.String("password", req.Password))

	tokens, err := uss.core.LoginUser(req.Password, req.Email)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("Login").Observe(float64(time.Since(then).Seconds()))

	return &pb.TokenPair{
		Refresh: tokens.RefreshToken,
		Access:  tokens.AccessToken,
	}, nil
}

func (uss *UserServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	then := time.Now()
	metrics.RequestTotal.WithLabelValues("RefreshToken").Inc()

	if req == nil {
		uss.logger.Error("empty request")

		return nil, ErrEmptyReq
	}

	uss.logger.Info("new refresh token request",
		zap.String("refresh_token", req.RefreshToken))

	accessToken, err := uss.core.Refresh(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("RefreshToken").Observe(float64(time.Since(then).Seconds()))

	return &pb.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (uss *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.TokenPair, error) {
	then := time.Now()
	metrics.RequestTotal.WithLabelValues("Register").Inc()

	if req == nil {
		uss.logger.Error("empty request")

		return nil, ErrEmptyReq
	}

	uss.logger.Info("new register request",
		zap.String("email", req.Email),
		zap.String("password", req.Password),
		zap.String("username", req.Username))

	tokens, err := uss.core.RegisterUser(req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}

	metrics.RequestDuration.WithLabelValues("Register").Observe(float64(time.Since(then).Seconds()))

	return &pb.TokenPair{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}, nil
}

func (uss *UserServiceServer) IncLike(ctx context.Context, req *pb.IncLikeRequest) (*emptypb.Empty, error) {
	if req == nil {
		uss.logger.Error("empty request")

		return &emptypb.Empty{}, ErrEmptyReq
	}

	uss.logger.Info("new inc like request",
		zap.String("id", req.UserId))

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.UserId))

		return &emptypb.Empty{}, ErrInvalidUUID
	}

	if err = uss.core.IncLike(uid); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (uss *UserServiceServer) DecLike(ctx context.Context, req *pb.DecLikeRequest) (*emptypb.Empty, error) {
	if req == nil {
		uss.logger.Error("empty request")

		return &emptypb.Empty{}, ErrEmptyReq
	}

	uss.logger.Info("new dec like request",
		zap.String("id", req.UserId))

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.UserId))

		return &emptypb.Empty{}, ErrInvalidUUID
	}

	if err = uss.core.DecLike(uid); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (uss *UserServiceServer) IncReview(ctx context.Context, req *pb.IncReviewRequest) (*emptypb.Empty, error) {
	if req == nil {
		uss.logger.Error("empty request")

		return &emptypb.Empty{}, ErrEmptyReq
	}

	uss.logger.Info("new inc like request",
		zap.String("id", req.UserId))

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.UserId))

		return &emptypb.Empty{}, ErrInvalidUUID
	}

	if err = uss.core.IncReview(uid); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (uss *UserServiceServer) DecReview(ctx context.Context, req *pb.DecReviewRequest) (*emptypb.Empty, error) {
	if req == nil {
		uss.logger.Error("empty request")

		return &emptypb.Empty{}, ErrEmptyReq
	}

	uss.logger.Info("new dec review request",
		zap.String("id", req.UserId))

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		uss.logger.Error("failed to parse uuid from request",
			zap.String("id", req.UserId))

		return &emptypb.Empty{}, ErrInvalidUUID
	}

	if err = uss.core.DecReview(uid); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
