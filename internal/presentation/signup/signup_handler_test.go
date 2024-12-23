package signup_test

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
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/presentation/signup"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func TestSignupHandler(t *testing.T) {
	var (
		userID          = idVO.NewUserIDForTest("user")
		userName        = "sato taro"
		userEmail       = "sato@example.com"
		userPassword    = "password"
		accessToken     = "token"
		invalidJSONBody = "invalid json"
		uri             = "/api/v1/signup"
		arg             = gomock.Any()
	)

	var happyRequestBody = signup.SignupRequest{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		prepare              func(mockSignupUC *appMock.MockISignupUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Positive: サインアップに成功する",
			requestBody: happyRequestBody,
			prepare: func(mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(arg, arg).Return(&authApp.SignupDTO{
					User: authApp.SignupUserDTO{
						ID:    userID.String(),
						Name:  userName,
						Email: userEmail,
					},
					AccessToken: accessToken,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: signup.SignupResponse{
				User: signup.SignupResponseBodyUser{
					ID:    userID.String(),
					Name:  userName,
					Email: userEmail,
				},
				AccessToken: accessToken,
			},
		},
		{
			caseName:     "Negative: リクエストボディが無効なJSONの場合、Bad Request を返す",
			requestBody:  invalidJSONBody,
			prepare:      func(mockSignupUC *appMock.MockISignupUsecase) {},
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
			requestBody:  signup.SignupRequest{},
			prepare:      func(mockSignupUC *appMock.MockISignupUsecase) {},
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
			caseName:    "Negative: ユーザーのメールアドレスが既に存在する場合、Conflict を返す",
			requestBody: happyRequestBody,
			prepare: func(mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(arg, arg).Return(nil, userDomain.ErrEmailAlreadyExists)
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
			caseName:    "Negative: 認証情報が既に存在する場合、Conflict を返す",
			requestBody: happyRequestBody,
			prepare: func(mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(arg, arg).Return(nil, authDomain.ErrAlreadyExists)
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
			caseName:    "Negative: 未知のエラーが発生した場合、Internal Server Error を返す",
			requestBody: happyRequestBody,
			prepare: func(mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(arg, arg).Return(nil, assert.AnError)
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
			tt.prepare(mockSignupUC)

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
