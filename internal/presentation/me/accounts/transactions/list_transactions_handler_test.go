package transactions_test

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
	transactionApp "github.com/u104rak1/pocgo/internal/application/transaction"
	"github.com/u104rak1/pocgo/internal/config"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/presentation/me/accounts/transactions"
	"github.com/u104rak1/pocgo/internal/server/response"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestListTransactionsHandler(t *testing.T) {
	var (
		mockAny        = gomock.Any()
		userID         = idVO.NewUserIDForTest("user")
		accountID      = idVO.NewAccountIDForTest("account")
		from           = "20240101"
		to             = "20241231"
		operationTypes = transactionDomain.Deposit + "," + transactionDomain.Withdrawal + "," + transactionDomain.Transfer
		sort           = "DESC"
		limit          = "10"
		page           = "1"
		transactionID  = idVO.NewTransactionIDForTest("transaction")
		transactionAt  = timer.GetFixedDateString()
		uri            = "/api/v1/me/accounts/" + accountID.String() + "/transactions"
	)

	tests := []struct {
		caseName             string
		requestQuery         string
		setupContext         func() context.Context
		prepare              func(mockListTransactionsUC *appMock.MockIListTransactionsUsecase)
		expectedCode         int
		expectedResponseBody interface{}
	}{
		{
			caseName:     "Positive: 取引一覧取得に成功する",
			requestQuery: "?from=" + from + "&to=" + to + "&operation_types=" + operationTypes + "&sort=" + sort + "&limit=" + limit + "&page=" + page,
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockListTransactionsUC *appMock.MockIListTransactionsUsecase) {
				mockListTransactionsUC.EXPECT().Run(mockAny, mockAny).Return(&transactionApp.ListTransactionsDTO{
					Total: 1,
					Transactions: []transactionApp.ListTransactionDTO{
						{
							ID:                transactionID.String(),
							AccountID:         accountID.String(),
							ReceiverAccountID: nil,
							OperationType:     transactionDomain.Deposit,
							Amount:            1000,
							Currency:          money.JPY,
							TransactionAt:     transactionAt,
						},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResponseBody: transactions.ListTransactionsResponse{
				Total: 1,
				Transactions: []transactions.ListTransactionsTransaction{
					{
						ID:            transactionID.String(),
						AccountID:     accountID.String(),
						OperationType: transactionDomain.Deposit,
						Amount:        1000,
						Currency:      money.JPY,
						TransactionAt: transactionAt,
					},
				},
			},
		},
		{
			caseName:     "Negative: クエリパラメータが無効な場合、Validation Failed を返す",
			requestQuery: "?from=invalid&to=invalid&operation_types=invalid&sort=invalid&limit=-1&page=-1",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare:      func(mockListTransactionsUC *appMock.MockIListTransactionsUsecase) {},
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
			caseName: "Negative: ユーザーIDがコンテキストに存在しない場合、Unauthorized を返す",
			setupContext: func() context.Context {
				return context.Background()
			},
			prepare:      func(mockListTransactionsUC *appMock.MockIListTransactionsUsecase) {},
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
			caseName: "Negative: 口座が存在しない場合、Not Found を返す",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockListTransactionsUC *appMock.MockIListTransactionsUsecase) {
				mockListTransactionsUC.EXPECT().Run(mockAny, mockAny).Return(nil, accountDomain.ErrNotFound)
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
			caseName: "Negative: 未知のエラーが発生した場合、Internal Server Error を返す",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), config.CtxUserIDKey(), userID.String())
				return ctx
			},
			prepare: func(mockListTransactionsUC *appMock.MockIListTransactionsUsecase) {
				mockListTransactionsUC.EXPECT().Run(mockAny, mockAny).Return(nil, assert.AnError)
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
			req := httptest.NewRequest(http.MethodGet, uri+tt.requestQuery, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("account_id")
			ctx.SetParamValues(accountID.String())
			ctx.SetRequest(req.WithContext(tt.setupContext()))

			mockListTransactionsUC := appMock.NewMockIListTransactionsUsecase(ctrl)
			tt.prepare(mockListTransactionsUC)

			h := transactions.NewListTransactionsHandler(mockListTransactionsUC)
			err := h.Run(ctx)

			if tt.expectedCode == http.StatusOK {
				assert.NoError(t, err)
				var resp transactions.ListTransactionsResponse
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
