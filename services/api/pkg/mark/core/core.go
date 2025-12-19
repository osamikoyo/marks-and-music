package core

import (
	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/music-and-marks/services/api/pkg/mark/handler"
)

type MarkCore struct{
	handler *handler.Handler
}

func SetupMarkCore()

func (m *MarkCore) RegisterHandler(e *echo.Echo) {

}