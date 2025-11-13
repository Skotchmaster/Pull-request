package auth

import (
	"net/http"
	"pull_request/internal/logging"

	httpserv "pull_request/internal/http"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) RequireLogin(logger logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := extractBearerToken(c.Request().Header.Get("Authorization"))
			if err != nil {
				logger.Errorf("invalid or missing token", "error", err)
				return c.JSON(http.StatusUnauthorized, httpserv.NewAPIError(httpserv.ErrCodeUnauthorized, httpserv.ErrMsgUnauthorized))
			}

			claims, err := m.parseToken(token)
			if err != nil {
				logger.Errorf("invalid or missing token", "error", err)
				return c.JSON(http.StatusUnauthorized, httpserv.NewAPIError(httpserv.ErrCodeUnauthorized, httpserv.ErrMsgUnauthorized))
			}

			user := &UserCtx{
				UserID: claims.Subject,
				Role:   claims.Role,
			}

			c.Set(userCtxKey, user)

			return next(c)
		}
	}
}

func (m *Middleware) RequireAdmin(logger logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := extractBearerToken(c.Request().Header.Get("Authorization"))
			if err != nil {
				logger.Errorf("invalid or missing token", "error", err)
				return c.JSON(http.StatusUnauthorized, httpserv.NewAPIError(httpserv.ErrCodeUnauthorized, httpserv.ErrMsgUnauthorized))
			}

			claims, err := m.parseToken(token)
			if err != nil {
				logger.Errorf("invalid or missing token", "error", err)
				return c.JSON(http.StatusUnauthorized, httpserv.NewAPIError(httpserv.ErrCodeUnauthorized, httpserv.ErrMsgUnauthorized))
			}

			if claims.Role != "admin" {
				logger.Infof("forbidden: non-admin user tried to access admin endpoint", "user_id", claims.Subject, "role", claims.Role)
				return c.JSON(http.StatusUnauthorized, httpserv.NewAPIError(httpserv.ErrCodeForbidden, httpserv.ErrMsgAdminRequired))
			}

			user := &UserCtx{
				UserID: claims.Subject,
				Role:   claims.Role,
			}

			c.Set(userCtxKey, user)

			return next(c)
		}
	}
}
