package handler

import (
	"github.com/osamikoyo/music-and-marks/services/api/pkg/music/client"
)

const DefaultPageSize = 30

type Handler struct {
	cc *client.MusicClient
}

func NewHandler(cc *client.MusicClient) *Handler {
	return &Handler{
		cc: cc,
	}
}
