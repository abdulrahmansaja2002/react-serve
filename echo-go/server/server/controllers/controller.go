package controllers

import (
	"echo-react-serve/server/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller interface {
	Handle(ctx echo.Context) error
}

type controller struct {
	s services.Service
}

func NewController(s services.Service) Controller {
	return &controller{s: s}
}

func (c *controller) Handle(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK!")
}
