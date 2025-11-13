package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	JWTSecret []byte
}

type Claims struct {
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type UserCtx struct {
	UserID string
	Role   string
}

const userCtxKey = "user"


func GetUser(c echo.Context) (*UserCtx, bool) {
	v := c.Get(userCtxKey)
	if v == nil {
		return nil, false
	}
	u, ok := v.(*UserCtx)
	return u, ok
}

func extractBearerToken(h string) (string, error) {
	if h == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization header")
	} 

	if parts[1] == "" {
		return "", errors.New("empty token")
	}

	return parts[1], nil
}

func(m *Middleware) parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return m.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

