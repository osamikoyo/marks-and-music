package handler

import "github.com/osamikoyo/music-and-marks/services/api/pkg/user/client"

const RefreshTokenCookieName = "music-and-marks-refresh"

type Handler struct {
	cc *client.UserClient
}

func NewHandler(cc *client.UserClient) *Handler {
	return &Handler{
		cc: cc,
	}
}
