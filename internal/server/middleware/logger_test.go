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
	"github.com/u104rak1/pocgo/internal/server/middleware"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func TestSetLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		handler    func(ctx echo.Context) error
		wantStatus int
		wantLevel  string
		wantMsg    string
	}{
		{
			name: "Positive: 正常なリクエストの場合、INFOレベルのログが出力される",
			path: "/200",
			handler: func(ctx echo.Context) error {
				return ctx.JSON(http.StatusOK, map[string]string{"message": "Success"})
			},
			wantStatus: http.StatusOK,
			wantLevel:  "INFO",
			wantMsg:    "request received",
		},
		{
			name: "Positive: バリデーションエラーの場合、WARNレベルのログが出力される",
			path: "/400",
			handler: func(ctx echo.Context) error {
				return response.BadRequest(ctx, assert.AnError)
			},
			wantStatus: http.StatusBadRequest,
			wantLevel:  "WARN",
			wantMsg:    assert.AnError.Error(),
		},
		{
			name: "Positive: クライアントエラーの場合、WARNレベルのログが出力される",
			path: "/400",
			handler: func(ctx echo.Context) error {
				return response.BadRequest(ctx, assert.AnError)
			},
			wantStatus: http.StatusBadRequest,
			wantLevel:  "WARN",
			wantMsg:    assert.AnError.Error(),
		},
		{
			name: "Positive: サーバーエラーの場合、ERRORレベルのログが出力される",
			path: "/500",
			handler: func(ctx echo.Context) error {
				return response.InternalServerError(ctx, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantLevel:  "ERROR",
			wantMsg:    assert.AnError.Error(),
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

			assert.Equal(t, tt.wantStatus, rec.Code)

			var logEntry map[string]interface{}
			err := json.NewDecoder(strings.NewReader(buf.String())).Decode(&logEntry)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantLevel, logEntry["level"])
			assert.Equal(t, tt.wantMsg, logEntry["msg"])
			assert.Equal(t, "GET", logEntry["method"])
			assert.Equal(t, tt.path, logEntry["uri"])
			assert.Equal(t, tt.wantStatus, int(logEntry["status"].(float64)))
			assert.Contains(t, logEntry, "user_agent")
			assert.Contains(t, logEntry, "client_ip")
			assert.Contains(t, logEntry, "request_id")
			assert.Contains(t, logEntry, "latency")
		})
	}
}
