package me_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/config"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/presentation/me"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
)

func TestReadMyProfileHandler(t *testing.T) {
	var (
		userID     = "01J9R7YPV1FH1V0PPKVSB5C8FW"
		userName   = "Sato Taro"
		userEmail  = "sato@example.com"
		unknownErr = errors.New("unknown error")
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
			expectedResponseBody: me.ReadMyProfileResponseBody{
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
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.UnauthorizedReason,
				Message: config.ErrUserIDMissing.Error(),
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
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.NotFoundReason,
				Message: userDomain.ErrNotFound.Error(),
			},
		},
		{
			caseName: "Unknown error occurs during profile retrieval.",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockReadUserUC *appMock.MockIReadUserUsecase) {
				mockReadUserUC.EXPECT().Run(ctx, userApp.ReadUserCommand{ID: userID}).Return(nil, unknownErr)
			},
			expectedCode: http.StatusInternalServerError,
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
			req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetRequest(req.WithContext(tt.setupContext()))

			mockReadUserUC := appMock.NewMockIReadUserUsecase(ctrl)
			tt.prepare(tt.setupContext(), mockReadUserUC)

			h := me.NewReadMyProfileHandler(mockReadUserUC)
			err := h.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusOK {
				var resp me.ReadMyProfileResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
			} else {
				var resp response.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
			}
		})
	}
}
