package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DeleteReview(c echo.Context) error {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)
	if err != nil{
		return c.String(http.StatusBadRequest, "faield convert id to int")
	}

	if err = h.cc.DeleteReview(c.Request().Context(), uint(id));err != nil{
		return c.String(http.StatusInternalServerError, "faield delete review " + err.Error())
	}

	return c.String(http.StatusOK, "deleted successfully")
}