package core

import (
	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
)

type Repository interface {
	GetAlbumByID(uid uuid.UUID) (*entity.Album, error)
	GetArtistByID(uid *uuid.UUID) (*entity.Artist, error)
	ReadAlbums(page_size, page_index int) ([]entity.Album, error)
	ReadArtist(page_size, page_index int) ([]entity.Artist, error)
}
