package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/services/user/entity"
)

func (h *Handler) Register(c echo.Context) error {
	var user entity.User

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "failed convert user")
	}

	tokens, err := h.cc.Register(c.Request().Context(), &user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed register user "+err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = RefreshTokenCookieName
	cookie.Value = tokens.RefreshToken
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, tokens)
}
