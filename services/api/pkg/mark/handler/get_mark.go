package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetMark(c echo.Context) error {
	releaseId := c.Param("releaseid")

	ctx := c.Request().Context()

	mark, err := h.cc.GetMark(ctx, releaseId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed get mark: "+err.Error())
	}

	return c.JSON(http.StatusOK, mark)
}
