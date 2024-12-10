package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/u104rak1/pocgo/internal/config"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func AuthorizationMiddleware(authService authDomain.IAuthenticationService, jwtSecretKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return response.Unauthorized(c, authDomain.ErrAuthorizationHeaderMissingOrInvalid)
			}

			accessToken := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err := authService.GetUserIDFromAccessToken(c.Request().Context(), accessToken, jwtSecretKey)
			if err != nil {
				return response.Unauthorized(c, err)
			}

			ctx := context.WithValue(c.Request().Context(), config.CtxUserIDKey(), userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
