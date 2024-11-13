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
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	"github.com/ucho456job/pocgo/internal/config"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/me/accounts"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
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
	)

	var requestBody = accounts.CreateAccountRequest{
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
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: response.ErrInvalidJSON.Error(),
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
		},
		{
			caseName:    "If the context does not have a user id, an authentication error will be returned.",
			requestBody: requestBody,
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(ctx context.Context, mockCreateAccountUC *appMock.MockICreateAccountUsecase) {},
			expectedCode: http.StatusUnauthorized,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.UnauthorizedReason,
				Message: config.ErrUserIDMissing.Error(),
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
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.NotFoundReason,
				Message: userDomain.ErrNotFound.Error(),
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
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: accountDomain.ErrLimitReached.Error(),
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/me/accounts", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetRequest(req.WithContext(tt.setupContext()))

			mockCreateAccountUC := appMock.NewMockICreateAccountUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockCreateAccountUC)

			h := accounts.NewCreateAccountHandler(mockCreateAccountUC)
			err = h.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp accounts.CreateAccountResponse
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
