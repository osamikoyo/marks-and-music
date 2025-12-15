package handler

import (
	"github.com/osamikoyo/music-and-marks/services/api/pkg/mark/client"
)

type Handler struct{
	cc *client.MarkClient
}

func NewHandler(cc *client.MarkClient) *Handler {
	return &Handler{
		cc: cc,
	}
}
