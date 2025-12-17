package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/services/mark/entity"
)

func (h *Handler) CreateReview(c echo.Context) error {
	var review entity.Review

	if err := c.Bind(&review); err != nil {
		return c.String(http.StatusBadRequest, "faield bind review")
	}

	if err := h.cc.CreateReview(c.Request().Context(), &review); err != nil {
		return c.String(http.StatusInternalServerError, "failed create review "+err.Error())
	}

	return c.String(http.StatusCreated, "created successfully")
}
