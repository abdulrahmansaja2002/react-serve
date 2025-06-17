package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	CSRFContextKey  = "csrf"
	CSRFTokenLookup = "header:X-CSRF-Token,form:_csrf,query:_csrf"
)

func csrfMiddlewares(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		csrf := middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: CSRFTokenLookup,
			ContextKey:  CSRFContextKey,
		})
		return csrf(next)(c)
	}
}
