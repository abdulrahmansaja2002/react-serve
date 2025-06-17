package routes

import (
	"echo-react-serve/helpers/middlewares"
	"echo-react-serve/web"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	middlewares.SetMainMiddleware(e)
	// e.Use(middlewares.CSRF) // un comment this line to enable CSRF protection
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	e.GET("/api/csrf", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{"csrf": c.Get(middlewares.CSRFContextKey)})
	})
	route(e)
	web.RegisterHandlers(e)
}
