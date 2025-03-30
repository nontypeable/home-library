package v1

import "github.com/labstack/echo/v4"

func (h *handler) UserRoutes(domain *echo.Group) {
	domain.POST("/sign-up", h.CreateUser)
	domain.POST("/sign-in", h.SignInUser)
}
