package transactions_test

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
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	transactionApp "github.com/u104rak1/pocgo/internal/application/transaction"
	"github.com/u104rak1/pocgo/internal/config"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/presentation/me/accounts/transactions"
	"github.com/u104rak1/pocgo/internal/server/response"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestExecuteTransactionHandler(t *testing.T) {
	var (
		accountID     = idVO.NewAccountIDForTest("account")
		userID        = idVO.NewUserIDForTest("user")
		password      = "1234"
		operationType = transactionDomain.Deposit
		amount        = 1000.0
		currency      = moneyVO.JPY
		transactionID = idVO.NewTransactionIDForTest("transaction")
		transactionAt = timer.GetFixedDateString()
		uri           = "/api/v1/me/accounts/" + accountID.String() + "/transactions"
		arg           = gomock.Any()
	)

	var happyRequestBody = transactions.ExecuteTransactionRequestBody{
		Password:      password,
		OperationType: operationType,
		Amount:        amount,
		Currency:      currency,
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		setupContext         func() context.Context
		prepare              func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Positive: 取引の実行に成功する",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(&transactionApp.ExecuteTransactionDTO{
					ID:            transactionID.String(),
					AccountID:     accountID.String(),
					OperationType: operationType,
					Amount:        amount,
					Currency:      currency,
					TransactionAt: transactionAt,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: transactions.ExecuteTransactionResponse{
				ID:            transactionID.String(),
				AccountID:     accountID.String(),
				OperationType: operationType,
				Amount:        amount,
				Currency:      currency,
				TransactionAt: transactionAt,
			},
		},
		{
			caseName:    "Negative: リクエストボディが無効なJSONの場合、Bad Request を返す",
			requestBody: "invalid json",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare:      func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {},
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
			caseName:    "Negative: バリデーションエラーが発生した場合、Bad Request を返す",
			requestBody: transactions.ExecuteTransactionRequestBody{},
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare:      func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {},
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
			caseName:    "Negative: コンテキストにユーザーIDがない場合、Unauthorized を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {},
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
			caseName:    "Negative: 通貨が異なる場合、Bad Request を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(nil, moneyVO.ErrDifferentCurrencyOperation)
			},
			expectedCode: http.StatusBadRequest,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLBadRequest,
				Title:    response.TitleBadRequest,
				Status:   http.StatusBadRequest,
				Detail:   moneyVO.ErrDifferentCurrencyOperation.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Negative: パスワードが不一致な場合、Forbidden を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(nil, accountDomain.ErrUnmatchedPassword)
			},
			expectedCode: http.StatusForbidden,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLForbidden,
				Title:    response.TitleForbidden,
				Status:   http.StatusForbidden,
				Detail:   accountDomain.ErrUnmatchedPassword.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Negative: 口座が見つからない場合、Not Found を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(nil, accountDomain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLNotFound,
				Title:    response.TitleNotFound,
				Status:   http.StatusNotFound,
				Detail:   accountDomain.ErrNotFound.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Negative: 受取口座が見つからない場合、Not Found を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(nil, accountDomain.ErrReceiverNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLNotFound,
				Title:    response.TitleNotFound,
				Status:   http.StatusNotFound,
				Detail:   accountDomain.ErrReceiverNotFound.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Negative: 残高が不足している場合、Unprocessable Entity を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(nil, moneyVO.ErrInsufficientBalance)
			},
			expectedCode: http.StatusUnprocessableEntity,
			expectedResponseBody: response.ProblemDetail{
				Type:     response.TypeURLUnprocessableEntity,
				Title:    response.TitleUnprocessableEntity,
				Status:   http.StatusUnprocessableEntity,
				Detail:   moneyVO.ErrInsufficientBalance.Error(),
				Instance: uri,
			},
		},
		{
			caseName:    "Negative: 未知のエラーが発生した場合、Internal Server Error を返す",
			requestBody: happyRequestBody,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockExecuteTransactionUC *appMock.MockIExecuteTransactionUsecase) {
				mockExecuteTransactionUC.EXPECT().Run(arg, arg).Return(nil, assert.AnError)
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
			ctx.SetParamNames("account_id")
			ctx.SetParamValues(accountID.String())
			ctx.SetRequest(req.WithContext(tt.setupContext()))

			mockExecuteTransactionUC := appMock.NewMockIExecuteTransactionUsecase(ctrl)
			tt.prepare(mockExecuteTransactionUC)

			h := transactions.NewExecuteTransactionHandler(mockExecuteTransactionUC)
			err = h.Run(ctx)

			if tt.expectedCode == http.StatusCreated {
				assert.NoError(t, err)
				var resp transactions.ExecuteTransactionResponse
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
