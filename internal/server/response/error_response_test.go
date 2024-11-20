package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/server/response"
)

func TestValidationFailed(t *testing.T) {
	e := echo.New()
	path := "/path/to/resource"
	req := httptest.NewRequest(http.MethodPost, path, nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	validationErrors := []response.ValidationError{
		{Field: "username", Message: "Username is required"},
		{Field: "email", Message: "Email format is invalid"},
	}

	err := response.ValidationFailed(ctx, validationErrors)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)

	resp, ok := he.Message.(response.ValidationProblemDetail)
	assert.True(t, ok)
	assert.Equal(t, response.ValidationProblemDetail{
		ProblemDetail: response.ProblemDetail{
			Type:     "https://example.com/probs/bad-request",
			Title:    "Validation Failed",
			Status:   http.StatusBadRequest,
			Detail:   "one or more validation errors occurred",
			Instance: path,
		},
		Errors: validationErrors,
	}, resp)
}

func TestErrorResponses(t *testing.T) {
	e := echo.New()
	path := "/path/to/resource"

	tests := []struct {
		caseName         string
		function         func(ctx echo.Context, err error) error
		expectedResponse response.ProblemDetail
	}{
		{
			caseName: "BadRequest",
			function: response.BadRequest,
			expectedResponse: response.ProblemDetail{
				Type:     "https://example.com/probs/bad-request",
				Title:    "Bad Request",
				Status:   http.StatusBadRequest,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Unauthorized",
			function: response.Unauthorized,
			expectedResponse: response.ProblemDetail{
				Type:     "https://example.com/probs/unauthorized",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Forbidden",
			function: response.Forbidden,
			expectedResponse: response.ProblemDetail{
				Type:     "https://example.com/probs/forbidden",
				Title:    "Forbidden",
				Status:   http.StatusForbidden,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "NotFound",
			function: response.NotFound,
			expectedResponse: response.ProblemDetail{
				Type:     "https://example.com/probs/not-found",
				Title:    "Not Found",
				Status:   http.StatusNotFound,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Conflict",
			function: response.Conflict,
			expectedResponse: response.ProblemDetail{
				Type:     "https://example.com/probs/conflict",
				Title:    "Conflict",
				Status:   http.StatusConflict,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "InternalServerError",
			function: response.InternalServerError,
			expectedResponse: response.ProblemDetail{
				Type:     "https://example.com/probs/internal-server-error",
				Title:    "Internal Server Error",
				Status:   http.StatusInternalServerError,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, path, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := tt.function(ctx, assert.AnError)
			he, ok := err.(*echo.HTTPError)
			assert.True(t, ok)

			resp, ok := he.Message.(response.ProblemDetail)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedResponse, resp)
		})
	}
}
