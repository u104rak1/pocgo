package signin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	"github.com/u104rak1/pocgo/internal/presentation/signin"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func TestSigninHandler(t *testing.T) {
	var (
		accessToken     = "token"
		invalidJSONBody = "invalid json"
		uri             = "/api/v1/signin"
		arg             = gomock.Any()
	)

	var happyRequestBody = signin.SigninRequest{
		Email:    "sato@example.com",
		Password: "password",
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		prepare              func(mockSigninUC *appMock.MockISigninUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Positive: サインインに成功する",
			requestBody: happyRequestBody,
			prepare: func(mockSigninUC *appMock.MockISigninUsecase) {
				mockSigninUC.EXPECT().Run(arg, arg).Return(&authApp.SigninDTO{
					AccessToken: accessToken,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: signin.SigninResponse{
				AccessToken: accessToken,
			},
		},
		{
			caseName:     "Negative: リクエストボディが無効なJSONの場合、Bad Request を返す",
			requestBody:  invalidJSONBody,
			prepare:      func(mockSigninUC *appMock.MockISigninUsecase) {},
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
			caseName:     "Negative: バリデーションが無効な場合、Validation Failed を返す",
			requestBody:  signin.SigninRequest{},
			prepare:      func(mockSigninUC *appMock.MockISigninUsecase) {},
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
			caseName:    "Negative: 認証に失敗した場合、Unauthorized を返す",
			requestBody: happyRequestBody,
			prepare: func(mockSigninUC *appMock.MockISigninUsecase) {
				mockSigninUC.EXPECT().Run(arg, arg).Return(nil, authDomain.ErrAuthenticationFailed)
			},
			expectedCode: http.StatusUnauthorized,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLUnauthorized,
				Title:    response.TitleUnauthorized,
				Status:   http.StatusUnauthorized,
				Detail:   authDomain.ErrAuthenticationFailed.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Negative: 未知のエラーが発生した場合、Internal Server Error を返す",
			requestBody: happyRequestBody,
			prepare: func(mockSigninUC *appMock.MockISigninUsecase) {
				mockSigninUC.EXPECT().Run(arg, arg).Return(nil, assert.AnError)
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
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSigninUC := appMock.NewMockISigninUsecase(ctrl)
			tt.prepare(mockSigninUC)

			h := signin.NewSigninHandler(mockSigninUC)
			err = h.Run(ctx)

			if tt.expectedCode == http.StatusCreated {
				assert.NoError(t, err)
				var resp signin.SigninResponse
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
