package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func TestValidationFailed(t *testing.T) {
	t.Run("Positive: ValidationFailed が正常に動作する", func(t *testing.T) {
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
				Type:     response.TypeURLValidationFailed,
				Title:    response.TitleValidationFailed,
				Status:   http.StatusBadRequest,
				Detail:   response.DetailValidationFailed,
				Instance: path,
			},
			Errors: validationErrors,
		}, resp)
	})
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
			caseName: "Positive: Bad Request が正常に動作する",
			function: response.BadRequest,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLBadRequest,
				Title:    response.TitleBadRequest,
				Status:   http.StatusBadRequest,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Positive: Unauthorized が正常に動作する",
			function: response.Unauthorized,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLUnauthorized,
				Title:    response.TitleUnauthorized,
				Status:   http.StatusUnauthorized,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Positive: Forbidden が正常に動作する",
			function: response.Forbidden,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLForbidden,
				Title:    response.TitleForbidden,
				Status:   http.StatusForbidden,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Positive: Not Found が正常に動作する",
			function: response.NotFound,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLNotFound,
				Title:    response.TitleNotFound,
				Status:   http.StatusNotFound,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Positive: Conflict が正常に動作する",
			function: response.Conflict,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLConflict,
				Title:    response.TitleConflict,
				Status:   http.StatusConflict,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Positive: Unprocessable Entity が正常に動作する",
			function: response.UnprocessableEntity,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLUnprocessableEntity,
				Title:    response.TitleUnprocessableEntity,
				Status:   http.StatusUnprocessableEntity,
				Detail:   assert.AnError.Error(),
				Instance: path,
			},
		},
		{
			caseName: "Positive: Internal Server Error が正常に動作する",
			function: response.InternalServerError,
			expectedResponse: response.ProblemDetail{
				Type:     response.TypeURLInternalServerError,
				Title:    response.TitleInternalServerError,
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
