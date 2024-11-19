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
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	validationErrors := []response.ValidationError{
		{Field: "username", Message: "Username is required"},
		{Field: "email", Message: "Email format is invalid"},
	}

	err := response.ValidationFailed(ctx, validationErrors)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, he.Code)

	resp, ok := he.Message.(response.ValidationErrorResponse)
	assert.True(t, ok)
	assert.Equal(t, response.ValidationFailedReason, resp.Reason)
	assert.Equal(t, validationErrors, resp.Errors)
}

func TestErrorResponses(t *testing.T) {
	e := echo.New()

	tests := []struct {
		caseName           string
		function           func(ctx echo.Context, err error) error
		expectedStatusCode int
		expectedResponse   response.ErrorResponse
	}{
		{
			caseName:           "BadRequest",
			function:           response.BadRequest,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: assert.AnError.Error(),
			},
		},
		{
			caseName:           "Unauthorized",
			function:           response.Unauthorized,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: response.ErrorResponse{
				Reason:  response.UnauthorizedReason,
				Message: assert.AnError.Error(),
			},
		},
		{
			caseName:           "Forbidden",
			function:           response.Forbidden,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse: response.ErrorResponse{
				Reason:  response.ForbiddenReason,
				Message: assert.AnError.Error(),
			},
		},
		{
			caseName:           "NotFound",
			function:           response.NotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: response.ErrorResponse{
				Reason:  response.NotFoundReason,
				Message: assert.AnError.Error(),
			},
		},
		{
			caseName:           "Conflict",
			function:           response.Conflict,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: assert.AnError.Error(),
			},
		},
		{
			caseName:           "InternalServerError",
			function:           response.InternalServerError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: response.ErrorResponse{
				Reason:  response.InternalServerErrorReason,
				Message: assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := tt.function(ctx, assert.AnError)
			he, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedStatusCode, he.Code)

			resp, ok := he.Message.(response.ErrorResponse)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedResponse, resp)
		})
	}
}
