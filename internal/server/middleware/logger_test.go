package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/server/middleware"
)

func TestSetLoggerMiddleware(t *testing.T) {
	e := echo.New()

	// ログの出力をキャプチャするためのバッファを用意
	var logBuffer bytes.Buffer
	writer := io.MultiWriter(&logBuffer)

	// EchoのLoggerをキャプチャするためのカスタムLoggerを設定
	e.Logger.SetOutput(writer)
	middleware.SetLoggerMiddleware(e)

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// 実際にハンドラーを呼び出してミドルウェアが処理されるか確認
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test response")
	}
	if assert.NoError(t, handler(ctx)) {
		e.ServeHTTP(rec, req)

		// レスポンスコードを確認
		assert.Equal(t, http.StatusOK, rec.Code)

		// ログの出力を確認
		logContent := logBuffer.String()
		assert.Contains(t, logContent, `"method":"GET"`)
		assert.Contains(t, logContent, `"uri":"/test"`)
		assert.Contains(t, logContent, `"status":200`)
		assert.Contains(t, logContent, `"latency"`)
	}
}
