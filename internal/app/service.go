package app

import (
	"github.com/labstack/echo/v4"
	userHTTPDelivery "home-library/internal/services/user/delivery/http/v1"
	userRepository "home-library/internal/services/user/repository"
	userUseCases "home-library/internal/services/user/usecases"
	"net/http"
)

func (app *App) startService() error {
	domain := app.echo.Group("/api/v1")

	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	var (
		userRepo        = userRepository.NewRepository(app.db)
		userUC          = userUseCases.NewUseCase(userRepo)
		userHTTPHandler = userHTTPDelivery.NewHandler(userUC)
	)
	userHTTPHandler.UserRoutes(domain)

	return nil
}
