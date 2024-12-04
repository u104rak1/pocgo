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
	authApp "github.com/u104raki/pocgo/internal/application/authentication"
	appMock "github.com/u104raki/pocgo/internal/application/mock"
	authDomain "github.com/u104raki/pocgo/internal/domain/authentication"
	userDomain "github.com/u104raki/pocgo/internal/domain/user"
	"github.com/u104raki/pocgo/internal/presentation/signup"
	"github.com/u104raki/pocgo/internal/server/response"
	"github.com/u104raki/pocgo/pkg/ulid"
)

func TestSignupHandler(t *testing.T) {
	var (
		userID             = ulid.GenerateStaticULID("user")
		userName           = "sato taro"
		userEmail          = "sato@example.com"
		userPassword       = "password"
		accessToken        = "token"
		invalidRequestBody = "invalid json"
		uri                = "/api/v1/signup"
	)

	var requestBody = signup.SignupRequest{
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
			expectedResponseBody: signup.SignupResponse{
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
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLBadRequest,
				Title:    response.TitleBadRequest,
				Status:   http.StatusBadRequest,
				Detail:   response.ErrInvalidJSON.Error(),
				Instance: uri,
			},
		},
		{
			caseName:     "Error occurs during signup when validation failed.",
			requestBody:  signup.SignupRequest{},
			prepare:      func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedResponseBody: response.ValidationProblemDetail{
				ProblemDetail: response.ProblemDetail{
					Type:     response.TypeURLValidationFailed,
					Title:    response.TitleValidationFailed,
					Status:   http.StatusBadRequest,
					Detail:   response.DetailValidationFailed,
					Instance: uri,
				},
			},
		},
		{
			caseName:    "Error occurs during signup when user email already exists.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, userDomain.ErrEmailAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLConflict,
				Title:    response.TitleConflict,
				Status:   http.StatusConflict,
				Detail:   userDomain.ErrEmailAlreadyExists.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Error occurs during signup when authentication already exists.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, authDomain.ErrAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLConflict,
				Title:    response.TitleConflict,
				Status:   http.StatusConflict,
				Detail:   authDomain.ErrAlreadyExists.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Unknown error occurs during signup.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLInternalServerError,
				Title:    response.TitleInternalServerError,
				Status:   http.StatusInternalServerError,
				Detail:   assert.AnError.Error(),
				Instance: uri,
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
			req := httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSignupUC := appMock.NewMockISignupUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockSignupUC)

			h := signup.NewSignupHandler(mockSignupUC)
			err := h.Run(ctx)

			if tt.expectedCode == http.StatusCreated {
				assert.NoError(t, err)
				var resp signup.SignupResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
			} else {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, he.Code)
				switch resp := he.Message.(type) {
				case response.ProblemDetail:
					assert.Equal(t, tt.expectedResponseBody, resp)
				case response.ValidationProblemDetail:
					expected := tt.expectedResponseBody.(response.ValidationProblemDetail)
					assert.Equal(t, expected.ProblemDetail, resp.ProblemDetail)
					assert.Greater(t, len(resp.Errors), 0)
				default:
					t.Errorf("unexpected response: %v", resp)
				}
			}
		})
	}
}
