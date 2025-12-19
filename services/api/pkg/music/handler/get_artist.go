package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetArtist(c echo.Context) error {
	id := c.Param("id")

	artist, err := h.cc.GetArtist(c.Request().Context(), id)
	if err != nil{
		return c.String(http.StatusInternalServerError, "failed get artist " + err.Error())
	}

	return c.JSON(http.StatusOK, artist)
}