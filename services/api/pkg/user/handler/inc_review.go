package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) IncReview(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()

	if err := h.cc.IncLike(ctx, id); err != nil {
		return c.String(http.StatusInternalServerError, "failed inc revied "+err.Error())
	}

	return c.String(http.StatusOK, "incremented successfully")
}
