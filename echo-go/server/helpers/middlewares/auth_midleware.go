package middlewares

import (
	"context"
	"log"
	"net/http"

	"echo-react-serve/helpers/auth"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return &auth.JwtCustomClaims{}
			},
			SigningKey: []byte(auth.SigningKey),
			ErrorHandler: func(c echo.Context, err error) error {
				log.Printf("JWT Middleware Error: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt token")
			},
			TokenLookup: "header:Authorization:Bearer , query:token, cookie:jwt", // need a space after 'Bearer'
			SuccessHandler: func(c echo.Context) {
				claims, err := auth.GetClaims(c.Get("user"))
				if err != nil {
					log.Printf("error getting claims: %v\n", err)
					return
				}
				ctx := context.WithValue(c.Request().Context(), auth.ContextKey, claims)
				ctx = context.WithValue(ctx, auth.ContextUserId, claims.UserId)
				c.SetRequest(c.Request().WithContext(ctx))
			},
		}
		log.Printf("Authorization header: %v", c.Request().Header.Get("Authorization"))
		return echojwt.WithConfig(config)(next)(c)
	}
}
