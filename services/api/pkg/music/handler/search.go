package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Search(c echo.Context) error {
	query := c.QueryParam("query")
	indexstr := c.QueryParam("index")
	index, err := strconv.Atoi(indexstr)
	if err != nil {
		return c.String(http.StatusBadRequest, "failed convert page index")
	}

	ctx := c.Request().Context()

	result, err := h.cc.Search(ctx, query, int32(index), DefaultPageSize)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed search "+err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
