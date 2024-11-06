package signup_test

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
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/signup"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestSignupHandler(t *testing.T) {
	var (
		userID             = ulid.GenerateStaticULID("user")
		userName           = "sato taro"
		userEmail          = "sato@example.com"
		userPassword       = "password"
		accessToken        = "token"
		invalidRequestBody = "invalid json"
	)

	var requestBody = signup.SignupRequestBody{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		prepare              func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Successful signup.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, authApp.SignupCommand{
					Name:     userName,
					Email:    userEmail,
					Password: userPassword,
				}).Return(&authApp.SignupDTO{
					User: authApp.SignupUserDTO{
						ID:    userID,
						Name:  userName,
						Email: userEmail,
					},
					AccessToken: accessToken,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: signup.SignupResponseBody{
				User: signup.SignupResponseBodyUser{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
				},
				AccessToken: accessToken,
			},
		},
		{
			caseName:     "Error occurs during signup when request body is invalid json.",
			requestBody:  invalidRequestBody,
			prepare:      func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: response.ErrInvalidJSON.Error(),
			},
		},
		{
			caseName:     "Error occurs during signup when validation failed.",
			requestBody:  signup.SignupRequestBody{},
			prepare:      func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName:    "Error occurs during signup when user email already exists.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, userDomain.ErrEmailAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: userDomain.ErrEmailAlreadyExists.Error(),
			},
		},
		{
			caseName:    "Error occurs during signup when authentication already exists.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, authDomain.ErrAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: authDomain.ErrAlreadyExists.Error(),
			},
		},
		{
			caseName:    "Unknown error occurs during signup.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, assert.AnError)
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
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSignupUC := appMock.NewMockISignupUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockSignupUC)

			h := signup.NewSignupHandler(mockSignupUC)
			err := h.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp signup.SignupResponseBody
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
