package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DecReview(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()

	if err := h.cc.DecReview(ctx, id); err != nil {
		return c.String(http.StatusInternalServerError, "fialed decrement review "+err.Error())
	}

	return c.String(http.StatusOK, "decrement successfully")
}
