package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ChangePassword(c echo.Context) error {
	id := c.Param("id")

	var req struct{
		Old string `json:"old"`
		New string `json:"new"`
	}

	if err := c.Bind(req);err != nil{
		return c.String(http.StatusBadRequest, "failed bind passwords")
	}

	if err := h.cc.ChangePassword(c.Request().Context(), id, req.Old, req.New);err != nil {
		return c.String(http.StatusBadRequest, "failed change password " + err.Error())
	}

	return c.String(http.StatusOK, "password changed successfully")
}