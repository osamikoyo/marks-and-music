package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DecLike(c echo.Context) error {
	idstr := c.Param("id")

	if err := h.cc.DecLike(c.Request().Context(), idstr);err != nil{
		return c.String(http.StatusInternalServerError, "failed dec like " + err.Error())
	}

	return c.String(http.StatusOK, "decremented successfully")
}