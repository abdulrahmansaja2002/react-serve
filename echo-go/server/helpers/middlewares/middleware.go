package middlewares

import "github.com/labstack/echo/v4"

// set group Jwr middleware
func SetGroupJwtMiddleware(g *echo.Group) {
	g.Use(jwtMiddleware)
}

// add another function to set middleware here
func SetMainMiddleware(e *echo.Echo) {
	setMainMiddlewares(e)
}

func CSRF(next echo.HandlerFunc) echo.HandlerFunc { return csrfMiddlewares(next) }
