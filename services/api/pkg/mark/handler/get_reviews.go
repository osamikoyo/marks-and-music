package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetReviews(c echo.Context) error {
	releaseid := c.Param("releaseid")

	ctx := c.Request().Context()

	reviews, err := h.cc.GetReviews(ctx, releaseid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed get reviews "+err.Error())
	}

	return c.JSON(http.StatusOK, reviews)
}
