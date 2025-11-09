package repository

import (
	"errors"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
	"gorm.io/gorm"
)

type Repository struct{
	logger *logger.Logger
	db *gorm.DB
}

var (
	ErrInternal     = errors.New("internal error")
	ErrAlreadyExist = errors.New("user already exist")
	ErrNotFound     = errors.New("user not found")
)


func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db: db,
		logger: logger,
	}
}

func (r *Repository) CreateSong(song *entity.Song) error {
	
}