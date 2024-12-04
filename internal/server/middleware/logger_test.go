package middleware_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/u104raki/pocgo/internal/server/middleware"
	"github.com/u104raki/pocgo/internal/server/response"
)

func TestSetLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		handler        func(ctx echo.Context) error
		expectedStatus int
		expectedLevel  string
		expectedMsg    string
	}{
		{
			name: "Output INFO level log when normal case.",
			path: "/200",
			handler: func(ctx echo.Context) error {
				return ctx.JSON(http.StatusOK, map[string]string{"message": "Success"})
			},
			expectedStatus: http.StatusOK,
			expectedLevel:  "INFO",
			expectedMsg:    "request received",
		},
		{
			name: "Output WARN level log when a validation error occurs.",
			path: "/400",
			handler: func(ctx echo.Context) error {
				return response.BadRequest(ctx, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			expectedLevel:  "WARN",
			expectedMsg:    assert.AnError.Error(),
		},
		{
			name: "Output WARN level log when an client error occurs.",
			path: "/400",
			handler: func(ctx echo.Context) error {
				return response.BadRequest(ctx, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			expectedLevel:  "WARN",
			expectedMsg:    assert.AnError.Error(),
		},
		{
			name: "Output ERROR level log when a server error occurs.",
			path: "/500",
			handler: func(ctx echo.Context) error {
				return response.InternalServerError(ctx, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedLevel:  "ERROR",
			expectedMsg:    assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			log.SetOutput(buf)

			e := echo.New()
			middleware.SetLoggerMiddleware(e)
			e.GET(tt.path, tt.handler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var logEntry map[string]interface{}
			err := json.NewDecoder(strings.NewReader(buf.String())).Decode(&logEntry)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedLevel, logEntry["level"])
			assert.Equal(t, tt.expectedMsg, logEntry["msg"])
			assert.Equal(t, "GET", logEntry["method"])
			assert.Equal(t, tt.path, logEntry["uri"])
			assert.Equal(t, tt.expectedStatus, int(logEntry["status"].(float64)))
			assert.Contains(t, logEntry, "user_agent")
			assert.Contains(t, logEntry, "client_ip")
			assert.Contains(t, logEntry, "request_id")
			assert.Contains(t, logEntry, "latency")
		})
	}
}
