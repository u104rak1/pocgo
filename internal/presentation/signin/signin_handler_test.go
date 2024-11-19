package signin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	"github.com/ucho456job/pocgo/internal/presentation/signin"
	"github.com/ucho456job/pocgo/internal/server/response"
)

func TestSigninHandler(t *testing.T) {
	var (
		accessToken        = "token"
		invalidRequestBody = "invalid json"
	)

	var requestBody = signin.SigninRequest{
		Email:    "sato@example.com",
		Password: "password",
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		prepare              func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Successful signin.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {
				mockSigninUC.EXPECT().Run(ctx, gomock.Any()).Return(&authApp.SigninDTO{
					AccessToken: accessToken,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: signin.SigninResponse{
				AccessToken: accessToken,
			},
		},
		{
			caseName:     "Error occurs during signin when request body is invalid json.",
			requestBody:  invalidRequestBody,
			prepare:      func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: response.ErrInvalidJSON.Error(),
			},
		},
		{
			caseName:     "Error occurs during signin when validation failed.",
			requestBody:  signin.SigninRequest{},
			prepare:      func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName:    "Error occurs during signin when authentication failed.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {
				mockSigninUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, authDomain.ErrAuthenticationFailed)
			},
			expectedCode: http.StatusUnauthorized,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.UnauthorizedReason,
				Message: authDomain.ErrAuthenticationFailed.Error(),
			},
		},
		{
			caseName:    "Unknown error occurs during signin.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {
				mockSigninUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.InternalServerErrorReason,
				Message: assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/signin", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSigninUC := appMock.NewMockISigninUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockSigninUC)

			h := signin.NewSigninHandler(mockSigninUC)
			err = h.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp signin.SigninResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
			} else if rec.Code == http.StatusBadRequest && tt.requestBody != invalidRequestBody {
				var resp response.ValidationErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, response.ValidationFailedReason, resp.Reason)
				assert.NotEmpty(t, resp.Errors)
			} else {
				var resp response.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
			}
		})
	}
}
