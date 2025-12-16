package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := h.cc.GetUser(c.Request().Context(), id)
	if err != nil{
		return c.String(http.StatusInternalServerError, "failed fetch user " + err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

