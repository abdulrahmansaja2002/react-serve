package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || cookie == nil {
			token := generateToken()
			signedToken := signValue(token)

			newCookie := &http.Cookie{
				Name:     "token",
				Value:    signedToken,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Secure:   false, // Set to true for HTTPS
				MaxAge:   15 * 60,
			}
			c.SetCookie(newCookie)
		}
		return next(c)
	}
}

func RequireSignedToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || cookie == nil {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "Invalid or missing token"})
		}
		_, valid := validateSignedValue(cookie.Value)
		if !valid {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "Invalid or missing token"})
		}
		return next(c)
	}
}
