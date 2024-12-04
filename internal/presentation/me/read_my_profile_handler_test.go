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
	appMock "github.com/u104raki/pocgo/internal/application/mock"
	userApp "github.com/u104raki/pocgo/internal/application/user"
	"github.com/u104raki/pocgo/internal/config"

	userDomain "github.com/u104raki/pocgo/internal/domain/user"
	"github.com/u104raki/pocgo/internal/presentation/me"
	"github.com/u104raki/pocgo/internal/server/response"
	"github.com/u104raki/pocgo/pkg/ulid"
)

func TestReadMyProfileHandler(t *testing.T) {
	var (
		userID    = ulid.GenerateStaticULID("user")
		userName  = "Sato Taro"
		userEmail = "sato@example.com"
		uri       = "/api/v1/me"
	)

	tests := []struct {
		caseName             string
		setupContext         func() context.Context
		prepare              func(ctx context.Context, mockReadUserUC *appMock.MockIReadUserUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName: "Successful profile retrieval.",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(ctx, userApp.ReadUserCommand{ID: userID}).Return(&userApp.ReadUserDTO{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResponseBody: me.ReadMyProfileResponse{
				ID:    userID,
				Name:  userName,
				Email: userEmail,
			},
		},
		{
			caseName: "Error occurs when user id is missing in context.",
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(ctx context.Context, mockReadUserUC *appMock.MockIReadUserUsecase) {},
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
			caseName: "Error occurs during profile retrieval because user not found.",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(ctx, userApp.ReadUserCommand{ID: userID}).Return(nil, userDomain.ErrNotFound)
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
			caseName: "Unknown error occurs during profile retrieval.",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(ctx, userApp.ReadUserCommand{ID: userID}).Return(nil, assert.AnError)
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
			tt.prepare(tt.setupContext(), mockReadUserUC)

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
