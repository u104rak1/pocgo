package me_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	userApp "github.com/u104rak1/pocgo/internal/application/user"
	"github.com/u104rak1/pocgo/internal/config"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/presentation/me"
	"github.com/u104rak1/pocgo/internal/server/response"
)

func TestReadMyProfileHandler(t *testing.T) {
	var (
		userID    = idVO.NewUserIDForTest("user")
		userName  = "Sato Taro"
		userEmail = "sato@example.com"
		uri       = "/api/v1/me"
		arg       = gomock.Any()
	)

	tests := []struct {
		caseName             string
		setupContext         func() context.Context
		prepare              func(mockReadUserUC *appMock.MockIReadUserUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName: "Positive: マイプロフィールの取得に成功する",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(arg, arg).Return(&userApp.ReadUserDTO{
					ID:    userID.String(),
					Name:  userName,
					Email: userEmail,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResponseBody: me.ReadMyProfileResponse{
				ID:    userID.String(),
				Name:  userName,
				Email: userEmail,
			},
		},
		{
			caseName: "Negative: ユーザーIDがコンテキストに存在しない場合、Unauthorized を返す",
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(mockReadUserUC *appMock.MockIReadUserUsecase) {},
			expectedCode: http.StatusUnauthorized,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLUnauthorized,
				Title:    response.TitleUnauthorized,
				Status:   http.StatusUnauthorized,
				Detail:   config.ErrUserIDMissing.Error(),
				Instance: uri,
			},
		},
		{
			caseName: "Negative: ユーザーが見つからない場合、Not Found を返す",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(arg, arg).Return(nil, userDomain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLNotFound,
				Title:    response.TitleNotFound,
				Status:   http.StatusNotFound,
				Detail:   userDomain.ErrNotFound.Error(),
				Instance: uri,
			},
		},
		{
			caseName: "Negative: 未知のエラーが発生した場合、Internal Server Error を返す",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(arg, arg).Return(nil, assert.AnError)
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
			req := httptest.NewRequest(http.MethodGet, uri, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetRequest(req.WithContext(tt.setupContext()))

			mockReadUserUC := appMock.NewMockIReadUserUsecase(ctrl)
			tt.prepare(mockReadUserUC)

			h := me.NewReadMyProfileHandler(mockReadUserUC)
			err := h.Run(ctx)

			if tt.expectedCode == http.StatusOK {
				assert.NoError(t, err)
				var resp me.ReadMyProfileResponse
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
