package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
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
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ValidationErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.ValidationFailedReason, resp.Reason)
	assert.Equal(t, validationErrors, resp.Errors)
}

func TestErrorResponses(t *testing.T) {
	e := echo.New()

	tests := []struct {
		caseName           string
		function           func(ctx echo.Context, err error) error
		errMessage         string
		expectedStatusCode int
		expectedResponse   response.ErrorResponse
	}{
		{
			caseName:           "BadRequest",
			function:           response.BadRequest,
			errMessage:         "bad request error",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: "bad request error",
			},
		},
		{
			caseName:           "Unauthorized",
			function:           response.Unauthorized,
			errMessage:         "unauthorized access",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: response.ErrorResponse{
				Reason:  response.UnauthorizedReason,
				Message: "unauthorized access",
			},
		},
		{
			caseName:           "Forbidden",
			function:           response.Forbidden,
			errMessage:         "forbidden action",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse: response.ErrorResponse{
				Reason:  response.ForbiddenReason,
				Message: "forbidden action",
			},
		},
		{
			caseName:           "NotFound",
			function:           response.NotFound,
			errMessage:         "not found",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: response.ErrorResponse{
				Reason:  response.NotFoundReason,
				Message: "not found",
			},
		},
		{
			caseName:           "Conflict",
			function:           response.Conflict,
			errMessage:         "conflict error",
			expectedStatusCode: http.StatusConflict,
			expectedResponse: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: "conflict error",
			},
		},
		{
			caseName:           "InternalServerError",
			function:           response.InternalServerError,
			errMessage:         "internal server error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: response.ErrorResponse{
				Reason:  response.InternalServerErrorReason,
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := tt.function(ctx, errors.New(tt.errMessage))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, rec.Code)

			var resp response.ErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, resp)
		})
	}
}
