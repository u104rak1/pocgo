package accounts_test

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
	accountApp "github.com/u104rak1/pocgo/internal/application/account"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	"github.com/u104rak1/pocgo/internal/config"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/presentation/me/accounts"
	"github.com/u104rak1/pocgo/internal/server/response"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestCreateAccountHandler(t *testing.T) {
	var (
		accountID          = ulid.GenerateStaticULID("account")
		userID             = ulid.GenerateStaticULID("user")
		name               = "For work"
		password           = "1234"
		currency           = money.JPY
		updatedAt          = timer.Now().String()
		invalidRequestBody = "invalid json"
		uri                = "/api/v1/me/accounts"
	)

	var requestBody = accounts.CreateAccountRequestBody{
		Name:     name,
		Password: password,
		Currency: currency,
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		setupContext         func() context.Context
		prepare              func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase)
		expectedCode         int
		expectedReason       string
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Successful account creation.",
			requestBody: requestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {
				mockCreateAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(&accountApp.CreateAccountDTO{
					ID:        accountID,
					UserID:    userID,
					Name:      name,
					Currency:  currency,
					UpdatedAt: updatedAt,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: accounts.CreateAccountResponse{
				ID:        accountID,
				Name:      name,
				Balance:   0,
				Currency:  currency,
				UpdatedAt: updatedAt,
			},
		},
		{
			caseName:    "An error occurs during account creation when request body is invalid json.",
			requestBody: invalidRequestBody,
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {},
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
			caseName:    "An error occurs during account creation when validation failed.",
			requestBody: accounts.CreateAccountRequest{},
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {},
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
			caseName:    "If the context does not have a user id, an authentication error will be returned.",
			requestBody: requestBody,
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {},
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
			caseName:    "If the user is not found, a not-found error will occur.",
			requestBody: requestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {
				mockCreateAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, userDomain.ErrNotFound)
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
			caseName:    "If account creation limit is reached, a conflict error will be returned.",
			requestBody: requestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {
				mockCreateAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, accountDomain.ErrLimitReached)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLConflict,
				Title:    response.TitleConflict,
				Status:   http.StatusConflict,
				Detail:   accountDomain.ErrLimitReached.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "If an unknown error occurs, an internal server error is returned.",
			requestBody: requestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID)
				return ctx
			},
			prepare: func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {
				mockCreateAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, assert.AnError)
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/me/accounts", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetRequest(req.WithContext(tt.setupContext()))

			mockCreateAccountUC := appMock.NewMockICreateAccountUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockCreateAccountUC)

			h := accounts.NewCreateAccountHandler(mockCreateAccountUC)
			err = h.Run(ctx)

			if tt.expectedCode == http.StatusCreated {
				assert.NoError(t, err)
				var resp accounts.CreateAccountResponse
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
