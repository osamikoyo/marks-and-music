package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ReadReleases(c echo.Context) error {
	indexstr := c.QueryParam("index")

	index, err := strconv.Atoi(indexstr)
	if err != nil {
		return c.String(http.StatusBadRequest, "failed convert index")
	}

	releases, err := h.cc.ReadReleases(c.Request().Context(), int32(index), DefaultPageSize)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed read releases "+err.Error())
	}

	return c.JSON(http.StatusOK, releases)
}
