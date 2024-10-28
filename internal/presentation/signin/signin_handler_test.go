package signin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/signin"
)

func TestSigninHandler(t *testing.T) {
	var (
		accessToken = "token"
		unknownErr  = errors.New("unknown error")
	)

	var requestBody = signin.SigninRequestBody{
		Email:    "sato@example.com",
		Password: "password",
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		prepare              func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase)
		expectedCode         int
		expectedReason       string
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
			expectedResponseBody: signin.SigninResponseBody{
				AccessToken: accessToken,
			},
		},
		{
			caseName:     "Error occurs during signin when request body is invalid json.",
			requestBody:  "invalid json",
			prepare:      func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: "code=400, message=Unmarshal type error: expected=signin.SigninRequestBody, got=string, field=, offset=14, internal=json: cannot unmarshal string into Go value of type signin.SigninRequestBody",
			},
		},
		{
			caseName:       "Error occurs during signin when validation failed.",
			requestBody:    signin.SigninRequestBody{},
			prepare:        func(ctx context.Context, mockSigninUC *appMock.MockISigninUsecase) {},
			expectedCode:   http.StatusBadRequest,
			expectedReason: response.ValidationFailedReason,
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
				mockSigninUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, unknownErr)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedReason: response.UnauthorizedReason,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.InternalServerErrorReason,
				Message: unknownErr.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/signin", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSigninUC := appMock.NewMockISigninUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockSigninUC)

			h := signin.NewSigninHandler(mockSigninUC)
			err := h.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp signin.SigninResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
			} else if rec.Code == http.StatusBadRequest && tt.expectedReason == response.ValidationFailedReason {
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
