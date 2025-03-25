package server

import (
	"github.com/labstack/echo/v4/middleware"
	"home-library/pkg/config"
)

func NewEchoServer(cfg *config.HTTPServerConfig) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Debug = cfg.Debug

	return e
}
