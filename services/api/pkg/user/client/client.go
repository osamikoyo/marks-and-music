package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/user/api/proto/gen/pb"
	"github.com/osamikoyo/music-and-marks/services/user/entity"
	"go.uber.org/zap"
)

var ErrNilInput = errors.New("input is nil")

type UserClient struct {
	cc     pb.UserServiceClient
	logger *logger.Logger
}

func NewUserClient(cc pb.UserServiceClient, logger *logger.Logger) *UserClient {
	return &UserClient{
		cc:     cc,
		logger: logger,
	}
}

func (u *UserClient) Register(ctx context.Context, user *entity.User) (*entity.TokenPair, error) {
	if user == nil {
		return nil, ErrNilInput
	}

	tokens, err := u.cc.Register(ctx, &pb.RegisterRequest{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})
	if err != nil {
		u.logger.Error("failed register",
			zap.Any("user", user),
			zap.Error(err))

		return nil, fmt.Errorf("failed register: %w", err)
	}

	return &entity.TokenPair{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}, nil
}

func (u *UserClient) Login(ctx context.Context, email, password string) (*entity.TokenPair, error) {
	if email == "" || password == "" {
		return nil, ErrNilInput
	}

	tokens, err := u.cc.Login(ctx, &pb.LoginRequest{Email: email, Password: password})
	if err != nil {
		u.logger.Error("failed login",
			zap.String("email", email),
			zap.String("password", password),
			zap.Error(err))

		return nil, fmt.Errorf("failed login: %w", err)
	}

	return &entity.TokenPair{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}, nil
}

func (u *UserClient) GetUser(ctx context.Context, id string) (*entity.User, error) {
	if id == "" {
		return nil, ErrNilInput
	}

	resp, err := u.cc.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		u.logger.Error("failed get user",
			zap.String("id", id),
			zap.Error(err))

		return nil, fmt.Errorf("failed get user: %w", err)
	}

	uid, err := uuid.Parse(resp.Id)
	if err != nil {
		u.logger.Error("failed parse id from resp",
			zap.String("id", resp.Id),
			zap.Error(err))

		return nil, fmt.Errorf("failed parse id: %w", err)
	}

	return &entity.User{
		ID:        uid,
		Username:  resp.Username,
		Email:     resp.Email,
		Reviews:   int(resp.Reviews),
		Likes:     int(resp.Likes),
		CreatedAt: resp.CreatedAt.AsTime(),
		UpdatedAt: resp.UpdatedAt.AsTime(),
	}, nil
}

func (u *UserClient) ChangePassword(ctx context.Context, id, currentPassword, newPassword string) error {
	if id == "" || currentPassword == "" || newPassword == "" {
		return ErrNilInput
	}

	_, err := u.cc.ChangePassword(ctx, &pb.ChangePasswordRequest{
		Id:              id,
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	})
	if err != nil {
		u.logger.Error("failed change password",
			zap.String("id", id),
			zap.Error(err))

		return fmt.Errorf("failed change password: %w", err)
	}

	return nil
}

func (u *UserClient) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrNilInput
	}

	_, err := u.cc.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		u.logger.Error("failed delete user",
			zap.String("id", id),
			zap.Error(err))

		return fmt.Errorf("failed delete user: %w", err)
	}

	return nil
}

func (u *UserClient) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", ErrNilInput
	}

	resp, err := u.cc.RefreshToken(ctx, &pb.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		u.logger.Error("failed refresh token",
			zap.Error(err))

		return "", fmt.Errorf("failed refresh token: %w", err)
	}

	return resp.AccessToken, nil
}

func (u *UserClient) IncLike(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrNilInput
	}

	_, err := u.cc.IncLike(ctx, &pb.IncLikeRequest{UserId: userID})
	if err != nil {
		u.logger.Error("failed inc like",
			zap.String("user_id", userID),
			zap.Error(err))

		return fmt.Errorf("failed inc like: %w", err)
	}

	return nil
}

func (u *UserClient) DecLike(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrNilInput
	}

	_, err := u.cc.DecLike(ctx, &pb.DecLikeRequest{UserId: userID})
	if err != nil {
		u.logger.Error("failed dec like",
			zap.String("user_id", userID),
			zap.Error(err))

		return fmt.Errorf("failed dec like: %w", err)
	}

	return nil
}

func (u *UserClient) IncReview(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrNilInput
	}

	_, err := u.cc.IncReview(ctx, &pb.IncReviewRequest{UserId: userID})
	if err != nil {
		u.logger.Error("failed inc review",
			zap.String("user_id", userID),
			zap.Error(err))

		return fmt.Errorf("failed inc review: %w", err)
	}

	return nil
}

func (u *UserClient) DecReview(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrNilInput
	}

	_, err := u.cc.DecReview(ctx, &pb.DecReviewRequest{UserId: userID})
	if err != nil {
		u.logger.Error("failed dec review",
			zap.String("user_id", userID),
			zap.Error(err))

		return fmt.Errorf("failed dec review: %w", err)
	}

	return nil
}
