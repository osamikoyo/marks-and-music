package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()

	if err := h.cc.DeleteUser(ctx, id);err != nil{
		return c.String(http.StatusInternalServerError, "failed delete user " + err.Error())
	}

	return c.String(http.StatusOK, "deleted successfully")
}