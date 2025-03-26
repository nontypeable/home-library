package app

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *App) startService() error {
	domain := app.echo.Group("/api/v1")

	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return nil
}
