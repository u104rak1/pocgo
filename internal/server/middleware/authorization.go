package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	"github.com/u104rak1/pocgo/internal/config"
	jwtServ "github.com/u104rak1/pocgo/internal/infrastructure/jwt"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func AuthorizationMiddleware(authService authApp.IJWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return response.Unauthorized(c, jwtServ.ErrInvalidAccessToken)
			}

			accessToken := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err := authService.GetUserIDFromAccessToken(accessToken)
			if err != nil {
				return response.Unauthorized(c, err)
			}

			ctx := context.WithValue(c.Request().Context(), config.CtxUserIDKey(), userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
