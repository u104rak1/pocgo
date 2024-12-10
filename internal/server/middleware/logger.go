package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/u104rak1/pocgo/internal/server/response"
	"golang.org/x/exp/slog"
)

func SetLoggerMiddleware(e *echo.Echo) {
	logger := slog.New(slog.NewJSONHandler(log.Writer(), nil))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			latency := time.Since(v.StartTime)

			var logLevel slog.Level
			switch {
			case v.Status >= 500:
				logLevel = slog.LevelError
			case v.Status >= 400:
				logLevel = slog.LevelWarn
			default:
				logLevel = slog.LevelInfo
			}

			var msg string
			if v.Error != nil {
				if httpError, ok := v.Error.(*echo.HTTPError); ok {
					switch m := httpError.Message.(type) {
					case response.ProblemDetail:
						msg = m.Detail
					case response.ValidationProblemDetail:
						msg = response.FormatValidationErrors(m.Errors)
					default:
						msg = fmt.Sprintf("%v", m)
					}
				} else {
					msg = v.Error.Error()
				}
			} else {
				msg = "request received"
			}

			logger.Log(c.Request().Context(), logLevel, msg,
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("user_agent", c.Request().UserAgent()),
				slog.String("client_ip", c.RealIP()),
				slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				slog.String("latency", latency.String()),
			)

			return nil
		},
	}))
}
