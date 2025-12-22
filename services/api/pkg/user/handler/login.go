package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "failed bind request")
	}

	tokens, err := h.cc.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed login "+err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = RefreshTokenCookieName
	cookie.Value = tokens.RefreshToken
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens)
}
