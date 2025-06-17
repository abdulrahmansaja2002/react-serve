package auth

import (
	"context"
	"log"
	"time"

	customErrors "echo-react-serve/constants/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserId int    `json:"user_id"`
	Admin  bool   `json:"admin"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	expiredOffset = time.Hour * 6
	SigningKey    = "secret"
	ContextKey    = contextKey("jwt_claims")
	ContextUserId = contextKey("user_id")
)

var SignMethod = jwt.SigningMethodHS256

func NewToken(userId int, admin bool, email string) string {
	claims := JwtCustomClaims{
		UserId: userId,
		Admin:  admin,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredOffset)),
		},
	}

	token := jwt.NewWithClaims(SignMethod, &claims)
	jwtToken, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		log.Printf("Error signing JWT token: %v", err)
	}
	return jwtToken
}

func GetClaims(user interface{}) (*JwtCustomClaims, error) {
	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, customErrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, customErrors.ErrFailedToParseClaims
	}
	return claims, nil
}

func GetClaimsFromContext(ctx context.Context) JwtCustomClaims {
	return ctx.Value(ContextKey).(JwtCustomClaims)
}
