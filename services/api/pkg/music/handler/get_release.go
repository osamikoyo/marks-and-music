package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetRelease(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()

	release, err := h.cc.GetRelease(ctx, id)
	if err != nil{
		return c.String(http.StatusInternalServerError, "faield get release " + err.Error())
	}

	return c.JSON(http.StatusOK, release)
}