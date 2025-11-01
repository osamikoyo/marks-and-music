package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/user/entity"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInternal     = errors.New("internal error")
	ErrAlreadyExist = errors.New("user already exist")
	ErrNotFound     = errors.New("user not found")
)

type Repository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *entity.User) error {
	r.logger.Info("creating user...",
		zap.Any("user", user))

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		r.logger.Error("failed to create user",
			zap.Any("user", user),
			zap.Error(err))

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExist
		}

		return ErrInternal
	}

	r.logger.Info("user was created successfully",
		zap.Any("user", user))

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, update *entity.User) error {
	r.logger.Info("updating user...",
		zap.Any("update", update))

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Save(update).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		r.logger.Error("failed update user",
			zap.Any("update", update),
			zap.Error(err))

		return ErrInternal
	}

	r.logger.Info("user was updated successfully")

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	r.logger.Info("deleting user...",
		zap.String("id", id.String()))

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		r.logger.Error("failed to delete user",
			zap.String("id", id.String()),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}

		return ErrInternal
	}

	r.logger.Info("user was deleted successfully",
		zap.String("id", id.String()))

	return nil
}

func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	r.logger.Info("fetching user...",
		zap.String("id", id.String()))

	var user entity.User

	err := r.db.WithContext(ctx).First(&user).Error

	if err != nil {
		r.logger.Error("failed to fetch user",
			zap.String("id", id.String()),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	r.logger.Error("user was successfully fetched",
		zap.Any("user", user))

	return &user, nil
}

func (r *Repository) CheckUser(ctx context.Context, email, password string) (string, error) {
	r.logger.Info("checking user...",
		zap.String("email", email))

	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("user not found",
				zap.String("email", email))
			return "", ErrNotFound
		}

		r.logger.Error("failed to query user",
			zap.String("email", email),
			zap.Error(err))
		return "", ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		r.logger.Warn("invalid password",
			zap.String("email", email))
		return "", ErrNotFound
	}

	r.logger.Info("user authenticated successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", email))

	return user.ID.String(), nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.logger.Info("fetching user by username...",
		zap.String("username", username))

	var user entity.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		r.logger.Error("failed to fetch user by username",
			zap.String("username", username),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	r.logger.Info("user was successfully fetched by username",
		zap.String("username", username),
		zap.Any("user", user))

	return &user, nil
}