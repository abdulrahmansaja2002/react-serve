package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const ngrokSkipBrowserWarningHeader = "ngrok-skip-browser-warning"
const xCsrfTokenHeader = "X-CSRF-Token"

func setMainMiddlewares(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: false,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, ngrokSkipBrowserWarningHeader, xCsrfTokenHeader},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
		}))
	e.Use(middleware.Recover())
}
