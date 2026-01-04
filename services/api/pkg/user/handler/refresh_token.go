package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie("music-and-marks-refresh")
	if err != nil {
		return c.String(http.StatusBadGateway, "not found refresh token")
	}

	token, err := h.cc.RefreshToken(ctx, cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed generate access token "+err.Error())
	}

	msg := struct {
		Token string `json:"access_token"`
	}{
		Token: token,
	}

	return c.JSON(http.StatusOK, msg)
}
